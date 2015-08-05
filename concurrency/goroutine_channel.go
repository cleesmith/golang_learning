package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	c := make(chan string)

	go sleepy("sleepy1: ", 1000, c)
	fmt.Printf("%q\n", <-c)

	go sleepy("sleepy2: ", 2000, c)
	fmt.Printf("%q\n", <-c)

	go sleepy("sleepy3: ", 500, c)
	fmt.Printf("%q\n", <-c)

  // guess at how long to wait, i.e. this is synchronous
  time.Sleep(time.Duration(3500) * time.Millisecond)

  elapsed := time.Since(start)
  fmt.Printf("main elapsed: %s\n", elapsed)

	// fmt.Printf("press any key to exit...\n")
	// var input string
	// fmt.Scanln(&input)
}

func sleepy(msg string, sleep_ms int, c chan string) {
  start := time.Now()
	sleep := time.Duration(sleep_ms) * time.Millisecond
  c <- fmt.Sprintf("%s working for %d milliseconds", msg, sleep)
  // note: the a message is sent to the channel before the sleep, if
  //       reversed each "go sleepy" call blocks and waits for the message:
  //       time.Sleep(sleep)
  //       c <- fmt.Sprintf("%s about to sleep for %d milliseconds", msg, sleep)
	time.Sleep(sleep)
  elapsed := time.Since(start)
  fmt.Printf("%s elapsed: %s\n", msg, elapsed)
}
