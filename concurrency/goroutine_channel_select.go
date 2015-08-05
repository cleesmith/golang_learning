package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	sleep_durations := []int{500, 1000, 2000, 7000, 8100}
	// sleep_durations := []int{8100, 1000, 2500, 500}
	c := make(chan string)
	defer close(c) // close channel when main exits
	for index, duration := range sleep_durations {
		go sleepy(fmt.Sprintf("sleepy%d: ", index+1), duration, c)
	}
	fmt.Printf("starting %d sleepys\n", len(sleep_durations))

	// timeout_chan := time.After(1 * time.Second)
	stimem := time.Now()
	for i, v := range sleep_durations {
		select {
		case msg := <-c:
			fmt.Println("received: ", msg)
		// case <-timeout_chan:
		case <-time.After(3 * time.Second):
			fmt.Println("*** TIMEOUT ***")
		}
		etimem := time.Since(stimem)
		fmt.Printf("\tfor range: i=%v v=%v elapsed=%v\n\n", i, v, etimem.Seconds())
	}

	elapsed := time.Since(start)
	fmt.Printf("... %d sleepys ran in: %e\n", len(sleep_durations), elapsed.Seconds())
}

func sleepy(msg string, sleep_ms int, yawn chan string) {
	start := time.Now()
	sleep := time.Duration(sleep_ms) * time.Millisecond
	time.Sleep(sleep) // some sleepy work
	yawn <- fmt.Sprintf("%s slept for %s", msg, sleep)
	elapsed := time.Since(start)
	fmt.Printf("\t%s finished in: %s\n", msg, elapsed)
}
