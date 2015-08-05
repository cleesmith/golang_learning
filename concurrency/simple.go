package main

import (
	"fmt"
)

func main() {
	start := make(chan bool)
	fmt.Printf("making 100 goroutines ... press any key to exit\n")
	for i := 0; i < 100; i++ {
		go worker(start, i)
	}
	close(start)
	// all workers start running now
	fmt.Printf("all workers start running now\n")

	var input string
	fmt.Scanln(&input)
}

func worker(start chan bool, id int) {
	<-start
	// do work
	fmt.Printf("worker %d working\n", id)
}
