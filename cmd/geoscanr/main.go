package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sethgrid/pester"
)

// URLSet was generated 2018-08-01 15:01:03 by tir on sol.
type URLSet struct {
	XMLName        xml.Name `xml:"urlset"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	URL            []struct {
		Text string `xml:",chardata"`
		Loc  struct {
			Text string `xml:",chardata"` // https://geoscan.nrcan.gc....
		} `xml:"loc"`
		Lastmod struct {
			Text string `xml:",chardata"` // 2010-01-25, 2010-02-02, 2...
		} `xml:"lastmod"`
	} `xml:"url"`
}

var (
	sitemap  = flag.String("sitemap", "https://geoscan.nrcan.gc.ca/googlesitemapGCxml.xml", "file or link to sitemap")
	cacheDir = flag.String("cachedir", filepath.Join(".", ".geoscanr"), "cache for page downloads")
)

func stringSum(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// fetch fetches a link or retrieves it from local cache.
func fetch(link string) ([]byte, error) {
	fn := filepath.Join(*cacheDir, stringSum(link))

	if _, err := os.Stat(fn); os.IsNotExist(err) {
		if _, err := os.Stat(*cacheDir); os.IsNotExist(err) {
			if err := os.MkdirAll(*cacheDir, 0755); err != nil {
				return nil, err
			}
		}
		resp, err := pester.Get(link)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return nil, err
		}
		f, err := ioutil.TempFile("", "geoscanr-")
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(f, resp.Body); err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
		os.Rename(f.Name(), fn)
	} else {
		log.Printf("cache hit %s", link)
	}

	return ioutil.ReadFile(fn)
}

func main() {
	flag.Parse()

	var r io.Reader

	if strings.HasPrefix(*sitemap, "http") {
		resp, err := pester.Get(*sitemap)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode >= 400 {
			log.Fatalf("failed to fetch sitemap: http %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		r = resp.Body
	} else {
		f, err := os.Open(*sitemap)
		if err != nil {
			log.Fatal(err)
		}
		r = f
	}

	dec := xml.NewDecoder(r)
	dec.Strict = false

	var set URLSet
	if err := dec.Decode(&set); err != nil {
		log.Fatal(err)
	}

	for _, u := range set.URL {
		data, err := fetch(u.Loc.Text)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[%d] %s", len(data), u.Loc.Text)
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}
		m := make(map[string]string)
		doc.Find("tr").Each(func(_ int, s *goquery.Selection) {
			k, v := strings.TrimSpace(s.Find("th").Text()), strings.TrimSpace(s.Find("td").Text())
			if k == "" {
				return
			}
			if k == "Download" {
				v = s.Find(`td > a`).AttrOr("href", "")
			}
			m[k] = v
		})
		b, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
