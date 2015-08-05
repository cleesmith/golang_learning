package main

import (
	"fmt"
	"time"
)

func main() {
	mainChannel := make(chan time.Duration)
	delays := []time.Duration{
		8 * time.Second,
		15 * time.Millisecond,
		800 * time.Millisecond,
		2 * time.Second,
	}
	for _, d := range delays {
		go func(d time.Duration) {
			time.Sleep(d)
			mainChannel <- d
		}(d)
	}

	for range delays {
		select {
		case v := <-mainChannel:
			fmt.Println(v)
		case <-time.After(time.Second):
			fmt.Println("*** TIMEOUT ***")
		}
	}

	return
}
