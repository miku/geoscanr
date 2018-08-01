geoscanr: cmd/geoscanr/main.go
	go build -o $@ $<

clean:
	rm -f geoscanr

