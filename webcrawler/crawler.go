package main

import (
	"fmt"
)

type link struct {
	url   string
	depth int
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// Crawl crawls the web using `fetcher`, processing the content found.
// Crawling starts at `url` and follows links up to the given `depth`.
func Crawl(url string, depth int, fetcher Fetcher) {
	done := make(chan int)
	links := make(chan link)
	go slurp(url, depth, fetcher, links, done)
	slurpers := 1
	for slurpers > 0 {
		select {
		case lnk := <-links:
			go slurp(lnk.url, lnk.depth, fetcher, links, done)
			slurpers++
		case <-done:
			slurpers--
		}
	}
}

// seen keeps tracks of urls encountered so far
var seen = make(map[string]bool)

// slurp uses Fetcher `f` to get and process content from the given `url`.
// It sends embedded links to `links` and signals once to `done` when finished.
// If `url` has been processed already or if `depth` is <= 0, slurp doesn't do
// anything and finishes immediately.
func slurp(url string, depth int, f Fetcher, links chan link, done chan int) {
	defer func() { done <- 1 }()
	if depth <= 0 || seen[url] {
		return
	}
	seen[url] = true
	body, urls, err := f.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		links <- link{u, depth - 1}
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
