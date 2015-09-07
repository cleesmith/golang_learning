// how to capture connections from this program:
//  terminal 1 -- while true; do lsof -i -P | grep -i "established";  done > conns.txt
//  terminal 2 -- time go run limit_concurrency.go
//  nano -c conns.txt
package main

import (
	"bufio"
	"fmt"
	"net/http"
	// "os"
	"regexp"
	"time"
)

// var servers = []string{
// 	"http://www.google.com/search?client=ubuntu&channel=fs&q=go+language&ie=utf-8&oe=utf-8",
// 	"http://golang.org/",
// 	"htp://cellipede.com/",
// 	"http://cleesmith.github.io/",
// 	"http://golang.org/doc/go_tutorial.html"}

var servers = []string{
	"http://cellipede.com/",
	"http://cleesmith.github.io/"}

var maxOpenRequests int = 100

var rpmRegexp = regexp.MustCompile("href=\"[^\"]+\"")

func main() {
	search := getSearchTerm()
	out := make(chan string)
	done := make(chan int)
	go printer(out)
	for _, url := range servers {
		go searchURL(search, url, out, done)
	}
	for i := 0; i < len(servers); i++ {
		fmt.Printf("i=%v\n", i)
		aDone := <-done
		fmt.Printf("aDone=%v\n", aDone)
	}
	if false {
		time.Sleep(5 * time.Second)
	}
}

func getSearchTerm() string {
	return "spud"
}

func die(message string) {
	fmt.Println(message)
	// os.Exit(1)
}

func printer(out chan string) {
	for {
		fmt.Println("printer: before out channel receive:")
		fmt.Println(<-out)
		fmt.Println("printer: after out channel receive")
	}
}

var requestSemaphore = make(chan int, maxOpenRequests) // Integer chanel with a maximum queue size

func searchURL(search string, url string, out chan string, done chan int) {
	requestSemaphore <- 1 // Block until put in the semaphore queue
	defer func() {
		<-requestSemaphore // Dequeue from the semaphore
	}()
	defer func() {
		done <- 1 // Signal that function is done
	}()
	response, err := http.Get(url)
	if err != nil {
		die("Could not read from " + url + ":" + err.Error())
		// panic(err)
		return
	}
	defer response.Body.Close()
	if err == nil {
		bufferedReader := bufio.NewReader(response.Body)
		err = searchAll(url, bufferedReader, url, out)
		// fmt.Printf("searchURL: url=%v\n", url)
	}
}

func searchAll(search string, reader *bufio.Reader, httpRoot string, out chan string) error {
	out <- "sent to channel: out -- from searchAll: search=\"" + search + "\""
	return nil
}
