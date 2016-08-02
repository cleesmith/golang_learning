package main

import (
	"fmt"
	"time"
)

func main() {
	startingTime := time.Now().UTC()
	// do something
	time.Sleep(1 * time.Second)
	endingTime := time.Now().UTC()

	var duration time.Duration = endingTime.Sub(startingTime)
	fmt.Printf("elapsed=%T=%v\n", duration, duration)
}
