// how to run:
// go run webscrape.go http://cleesmith.github.io/
// go run webscrape.go http://cellipede.com/ http://cleesmith.github.io/
package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	_ "strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	seedUrls := os.Args[1:]

	foundUrls := make(map[string]bool)
	foundUAs := make(map[string]bool)

	chLinks := make(chan string)
	chUAs := make(chan string)
	chFin := make(chan bool)

	for _, url := range seedUrls {
		go scrape(url, chLinks, chUAs, chFin)
	}

	// gather results via channels
	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chLinks:
			foundUrls[url] = true
		case url := <-chUAs:
			foundUAs[url] = true
		case <-chFin:
			c++
		}
	}
	close(chLinks)
	fmt.Println("\nFound", len(foundUrls), "unique urls:")
	for url, _ := range foundUrls {
		fmt.Printf("\t%v\n", url)
	}
	close(chUAs)
	fmt.Println("\nFound", len(foundUAs), "UA's:")
	for ua, _ := range foundUAs {
		fmt.Printf("\t%v\n", ua)
	}
}

func scrape(url string, chL chan string, chU chan string, chDone chan bool) {
	fmt.Println("scrape: \"" + url + "\"")
	// resp, err := http.Get(url)
	// ... or with a timeout:
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)

	// ensure we send/notify to indicate we're done when this func (crawl) exits:
	defer func() { chDone <- true }()

	if err != nil {
		fmt.Println("ERROR: Failed to scrape \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // ensure resp.Body is closed when this func (crawl) exits

	// extract all http links from a webpage
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		fmt.Printf("z.Next=%T=%+v\n", tt, tt)

		switch tt {
		case html.ErrorToken:
			return // document processed so we're done

		case html.StartTagToken:
			t := z.Token()
			if t.Data == "script" {
				z.Next()
				tn := z.Token()
				// fmt.Printf("\ttn=%T=%v\n", tn, tn)
				// fmt.Printf("\ttn.Data=%T=%v\n", tn.Data, tn.Data)
				// fmt.Printf("\ttn.Attr=%T=%+v\n", tn.Attr, tn.Attr)
				// for k, v := range tn.Attr {
				// 	fmt.Printf("\t\ttn.Attr: k=%T=%+v v=%T=%+v\n", k, k, v, v)
				// }

				// yoink "UA-" values:
				// re, err := regexp.Compile("?") // test error
				re, err := regexp.Compile("(?:'|\")(?P<scrape>UA-.*?)(?:'|\")")
				if err != nil {
					continue
				}
				// m := re.MatchString(string(page)) // boolean
				uas := re.FindAllString(string(tn.Data), -1)
				if len(uas) > 0 {
					for _, ua := range uas {
						fmt.Printf(">>>>> ua=%T=%+v <<<<<\n", ua, ua)
						chU <- ua
					}
				}
				continue
			}

			// check if the token is an "<a ..." tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// make sure the url begins with http:
			// hasProto := strings.Index(url, "http") == 0
			// if hasProto {
			// 	ch <- url
			// }

			// to capture relative links as well:
			chL <- url
		}
	}
}

func getHref(t html.Token) (ok bool, href string) {
	// iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	// a "bare/empty" return will return the vars (ok, href) as defined in func getHref
	return
}
