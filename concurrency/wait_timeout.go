package main

import (
	"fmt"
	"log"
	"time"
)

// a function that will block and only return when all connections have ended:
func waitForConnections() chan string {
	c := make(chan string)
	go func() {
		time.Sleep(2 * time.Second) // fake work
		c <- fmt.Sprintf("slept for 2 seconds")
	}()
	return c
}

func main() {
	select {
	// wait for all connections to end
	case msg1 := <-waitForConnections():
		log.Println("received: ", msg1)
		log.Println("Waited for connections to end...")
	// don't want to wait indefinitely, so timeout after seconds elapsed:
	case <-time.After(1 * time.Second):
		log.Println("Timeout")
	}
}
