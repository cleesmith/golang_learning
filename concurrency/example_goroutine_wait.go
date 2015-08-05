package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"http://www.cnn.com",
		"http://espn.go.com/",
		"http://grantland.com",
		"http://www.newyorker.com/",
	}

	startIt := time.Now()
	without_goroutine(urls)
	endIt := time.Since(startIt)
	fmt.Printf("without goroutine elapsed time is %v\n\n", endIt.Seconds())

	startIt = time.Now()
	with_goroutine(urls)
	endIt = time.Since(startIt)
	fmt.Printf("with goroutine elapsed time is %v\n\n", endIt.Seconds())
}

func without_goroutine(urls []string) {
	for _, url := range urls {
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %s\n", url, response.Status)
	}
}

func with_goroutine(urls []string) {
	var wg sync.WaitGroup
	// number of concurrent goroutines to wait for:
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			response, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s: %s\n", url, response.Status)
		}(url)
	}
	wg.Wait() // waits until the url checks complete
}
