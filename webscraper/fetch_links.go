// usage:
// go run fetch_links.go --url http://cleesmith.github.io/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var (
	queryURL = flag.String("url", "", "URL to query")
)

type AnchorCollector struct {
	Base  string
	base  *url.URL
	hrefs map[string]struct{}
}

func (c *AnchorCollector) Selection() func(int, *goquery.Selection) {
	return func(_ int, s *goquery.Selection) {
		u, exists := s.Attr("href")
		if !exists {
			return
		}

		if c.base == nil {
			var err error
			if c.base, err = url.Parse(c.Base); err != nil {
				return
			}
		}

		urlobj, err := c.base.Parse(u)
		if err != nil {
			return
		}

		if urlobj.Scheme != "http" && urlobj.Scheme != "https" {
			return
		}

		if c.hrefs == nil {
			c.hrefs = make(map[string]struct{})
		}

		var x struct{}
		c.hrefs[urlobj.String()] = x
	}
}

func (c *AnchorCollector) Hrefs() []string {
	var res []string
	for u, _ := range c.hrefs {
		res = append(res, u)
	}

	return res
}

func main() {
	flag.Parse()
	if len(*queryURL) == 0 {
		log.Fatal("URL not set")
	}

	doc, err := goquery.NewDocument(*queryURL)
	if err != nil {
		log.Fatal(err)
	}

	c := AnchorCollector{Base: *queryURL}
	doc.Find("a").Each(c.Selection())
	for _, urlstr := range c.Hrefs() {
		fmt.Println(urlstr)
	}
}
