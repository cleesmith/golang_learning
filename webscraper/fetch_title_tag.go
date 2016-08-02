package main

import (
	"encoding/json"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Aug 2016: why use goquery ?
//           coz there are lots of malformed html on the internet,
//           and goquery is similar to using jQuery which makes
//           dealing with html (the dom) a lot simpler

type ScrapeResult struct {
	// note: names must be capitalized to be exported
	Url   string `json:"url"`
	Title string `json:"title"`
}

func main() {
	max_loops := 5
	// overview:
	// launch n goroutines
	//    each goroutine:
	//      gets html for URL
	//      scrapes title tag
	//      send title tag result via channel
	chResult := make(chan []byte)
	chDone := make(chan bool)

	for i := 0; i < max_loops; i++ {
		go scrape("http://cleesmith.github.io/", chResult, chDone)
	}

	// gather results via channels, similar to polling an SQS queue
	for c := 0; c < max_loops; {
		select {
		case json_result := <-chResult:
			sr := ScrapeResult{}
			if err := json.Unmarshal(json_result, &sr); err != nil {
				fmt.Printf("err=%v\n", err)
			}
			fmt.Printf("sr=%T=%+v\nsr.Url=%v\nsr.Title=%v\n", sr, sr, sr.Url, sr.Title)
		case <-chDone:
			c++ // too lame?
		}
	}
	close(chResult)
}

func scrape(aUrl string, result chan []byte, done chan bool) {
	// ensure we send/notify to indicate we're done when this func exits:
	defer func() { done <- true }()

	// get html page/doc via URL:
	doc, err := goquery.NewDocument(aUrl)
	if err != nil {
		err_msg := "ERROR: Failed to scrape \"" + aUrl + "\""
		j, json_err := json.Marshal(ScrapeResult{aUrl, err_msg})
		if json_err != nil {
			fmt.Println(json_err)
		}
		result <- j
		return
	}

	// only get the title tag from head (not in the body):
	title_tag_text := ""
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		title_tag_text = s.Find("title").Text()
	})

	var sr1 ScrapeResult
	sr1 = ScrapeResult{Url: aUrl, Title: title_tag_text}
	// sr1 = ScrapeResult{aUrl, title_tag_text}
	srj, jerr := json.Marshal(sr1)
	if jerr != nil {
		fmt.Println(jerr)
	}

	result <- srj
}
