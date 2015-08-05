package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c := boring("boring!")
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Printf("You're boring; I'm leaving.\n")
}

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() { // launch goroutine from inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // return the channel to the caller
}
