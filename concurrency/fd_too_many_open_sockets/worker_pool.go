// see: http://burke.libbey.me/conserving-file-descriptors-in-go/
// moderate File Descriptor consumption in two simple ways:
//  1. worker pool
//  2. semaphore
// the following is the worker pool strategy:
package main

import (
	"fmt"
	// "time"
)

const nWorkers = 2

var resp = make(chan string)

func main() {
	fmt.Printf("main: nWorkers=%v\n", nWorkers)
	up := make(chan string)
	crawl(up)
	for w := 1; w <= 4; w++ {
		ws := fmt.Sprintf("%v", w) + ".com"
		up <- ws
	}
	// for w := 1; w <= 0; w++ {
	// 	var aresp string
	// 	aresp = <-resp
	// 	fmt.Printf("main: aresp=%v\n", aresp)
	// }

	// time.Sleep(1 * time.Second)
}

func fetcher(url string) {
	fmt.Printf("fetcher: url=%v\n", url)
	// resp <- fmt.Sprintf("fetcher: url=%v\n", url)
	// fmt.Printf("fetcher: resp=%v\n", resp)
}

func crawl(urlProducer chan string) {
	fmt.Printf("crawl: urlProducer=%T\n", urlProducer)
	for i := 0; i < nWorkers; i++ {
		fmt.Printf("\ti=%v\n", i)
		go func() {
			for {
				fmt.Println("\tcalling fetcher")
				fetcher(<-urlProducer)
			}
		}()
	}
}
