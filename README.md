# geoscanr

A custom fetcher for geoscan data.

* [x] fetch sitemap (https://geoscan.nrcan.gc.ca/googlesitemapGCxml.xml)
* [x] fetch all links, extract data
* [x] cache html locally
* [ ] fetch linked assets, maybe
* [ ] create complete set in a single file

We do not need incremental updates, because we can use RSS for that.

```
$ geoscanr -h
  -cachedir string
        cache for page downloads (default ".geoscanr")
  -q    suppress logging output
  -sitemap string
        file or link to sitemap (default "https://geoscan.nrcan.gc.ca/googlesitemapGCxml.xml")

$ geoscanr | head -1 | jq .
{
  "Area": "Point Alexander; Renfrew County; Pontiac County; Nipissing District",
  "Author": "Canada Surveys and Mapping Branch, Topographical Survey / Direction canadienne des levés et de la cartographie, levés topographiques",
  "Document": "serial",
  "Download": "https://geoscan.nrcan.gc.ca/starweb/geoscan/servlet.starweb?path=geoscan/downloade.web&search1=R=123354",
  "Edition": "provisional/provisoire, black & white / noir et blanc",
  "File format": "pdf (Adobe® Reader®); JPEG2000",
  "GEOSCAN ID": "123354",
  "Illustrations": "location maps",
  "Image": "",
  "Lang.": "English",
  "Lat/Long WENS": "-77.7500  -77.5000   46.2500   46.0000",
  "Links": null,
  "Map Info.": "topographic, 1:50,000",
  "Maps": "1 map",
  "Media": "paper; on-line; digital",
  "NTS": "31K/04NE; 31K/04SE",
  "Province": "Quebec; Ontario",
  "Related": [
    "This publication supercedes Canada Bureau of Geology and\nTopography, Topographical Survey; (1942). ...""
  ],
  "Released": "1956 01 01",
  "Source": "Geological Survey of Canada, \"A\" Series Map 701A,   1946,  1 sheet, https://doi.org/10.4095/123354",
  "Subjects": "miscellaneous; topography; triangulation stations; hydrography; streams; marshes; reefs",
  "Title": "Point Alexander, Nipissing District, Renfrew and Pontiac counties, Ontario and Québec",
  "Year": "1946"
}
```
