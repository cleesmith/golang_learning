package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://cellipede.com:423/",
		"http://cellipede.com:4235/spud",
		"http://cellipede.com:4235/",
	}
	for _, url := range urls {
		// increment the WaitGroup counter:
		wg.Add(1)
		// launch a goroutine to fetch the URL:
		go func(url string) {
			// decrement the counter when the goroutine completes:
			defer wg.Done()
			// fetch the URL:
			// resp, err := http.Get(url)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
        return
    	}
    	// close connection to http server to avoid running out of sockets/file descriptors:
    	req.Close = true
			resp, err := http.DefaultClient.Do(req)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			req.Close = true
			fmt.Printf("\nurl=%s response:\n%v\n'___'\n", url, resp)
			fmt.Printf("resp.StatusCode=%d\n", resp.StatusCode)
		}(url) // this is the url passed in to the anonymous goroutine
	}
	// wait for all HTTP fetches to complete:
	wg.Wait()
	fmt.Println("all done!")
}
