package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Aug 2016: why use goquery ?
//           coz there are lots of malformed html on the internet,
//           and goquery is similar to using jQuery which makes
//           dealing with html (the dom) a lot simpler

type ScrapeResult struct {
	// note: names must be capitalized to be exported
	Url     string        `json:"url"`
	Title   string        `json:"title"`
	Elapsed time.Duration `json:"elapsed"`
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
		// go scrape("http://cleesmith.github.io/", chResult, chDone)
		go scrape("http://rctheatres.com/", chResult, chDone)
	}

	// since we know how many url's we are scraping (max_loops), we
	// know how many results to expect ... so let's just loop and
	// gather the results via channels, which is similar to polling
	// an SQS queue that we coded in boto3_test/fetch_title_via_sqs2.py
	for c := 0; c < max_loops; {
		select {
		case json_result := <-chResult:
			sr := ScrapeResult{}
			if err := json.Unmarshal(json_result, &sr); err != nil {
				fmt.Printf("err=%v\n", err)
			}
			// fmt.Printf("sr=%T=%+v\nsr.Url=%v\nsr.Title=%v\n", sr, sr, sr.Url, sr.Title)
			fmt.Printf("sr=%T=%+v\n", sr, sr)
		case <-chDone:
			c++ // too lame?
		}
	}
	close(chResult)
}

func scrape(aUrl string, result chan []byte, done chan bool) {
	startingTime := time.Now().UTC()
	// ensure we send/notify to indicate we're done when this func exits:
	defer func() { done <- true }()

	// get html page/doc via URL:
	doc, err := goquery.NewDocument(aUrl)
	if err != nil {
		err_msg := "ERROR: Failed to scrape \"" + aUrl + "\""
		j, json_err := json.Marshal(ScrapeResult{aUrl, err_msg, 0})
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

	endingTime := time.Now().UTC()
	var duration time.Duration = endingTime.Sub(startingTime)

	var sr1 ScrapeResult
	sr1 = ScrapeResult{Url: aUrl, Title: title_tag_text, Elapsed: duration}
	// sr1 = ScrapeResult{aUrl, title_tag_text}
	srj, jerr := json.Marshal(sr1)
	if jerr != nil {
		fmt.Println(jerr)
	}

	result <- srj
}
