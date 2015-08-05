package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool, 1)

	go func() {
		select {
		case m := <-c:
			fmt.Printf("m=%v\n", handle(m))
		case <-time.After(3 * time.Second):
			fmt.Println("timed out")
		}
	}()

	time.Sleep(2 * time.Second)
}

func handle(m bool) chan string {
	c := make(chan string)
	c <- "hi!"
	return c
}
