package main

import (
	"fmt"
	"time"
)

func main() {
	max_time := time.Duration(10000) * time.Nanosecond
	max_time = time.Duration(1) * time.Second
	atimeout := time.After(max_time) // creates a channel

	numbers := make(chan int)

	go func() {
		for n := 0; ; {
			numbers <- n // send to channel
			n++
		}
	}()

readChannel:
	for {
		select {
		case <-atimeout:
			fmt.Printf("... all done after %s\n", max_time)
			break readChannel
		case num := <-numbers:
			fmt.Println(num)
		}
	}

	fmt.Println("all done!")
}
