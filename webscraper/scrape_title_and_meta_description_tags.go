package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Aug 2016: why use goquery ?
//           coz there are lots of malformed html on the internet,
//           and goquery is similar to using jQuery which makes
//           dealing with html (the dom) a lot simpler

func main() {
	// doc via URL:
	// doc, err := goquery.NewDocument("http://cleesmith.github.io/")
	// doc via string, you know, for testing:
	r := strings.NewReader("<html><head><tItLe>Molly, a dog</TITLE></head><body></body></html>")
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic("wtf?")
	}

	// only get the title tag from head (not in the body):
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		pageTitle := s.Find("title").Text()
		fmt.Println(pageTitle)
	})

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
			description, _ := s.Attr("content")
			fmt.Println(description)
		}
	})
}
