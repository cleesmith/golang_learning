// how to run:
// go run webscrape.go http://cleesmith.github.io/
// go run webscrape.go http://cellipede.com:4235/ http://cleesmith.github.io/
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool) {
	fmt.Println("fetching: \"" + url + "\"")
	// resp, err := http.Get(url)
	// ... or with a timeout:
	// timeout := time.Duration(1 * time.Millisecond) // cause timeout
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)

	// notify that we're done when "crawl" function exits:
	defer func() { chFinished <- true }()

	if err != nil {
		fmt.Println("ERROR: Failed to scrape \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// make sure the url begins with http:
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}
	}
}

func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]
	chUrls := make(chan string)
	chFinished := make(chan bool)
	// get each page (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}
	// subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}
	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")
	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}
	close(chUrls)
}
