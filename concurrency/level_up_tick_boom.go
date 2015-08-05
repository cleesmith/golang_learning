package main

import (
	"fmt"
	"time"
)

func main() {
	stopChan := make(chan bool)

	go func() {
		time.Sleep(4100 * time.Millisecond)
		stopChan <- true
	}()

	timer := time.NewTimer(time.Second)

LOOP:
	for {
		select {
		case <-timer.C:
			// some polling code here
			fmt.Println("tick")
			timer.Reset(time.Second)
		case <-stopChan:
			fmt.Println("boom")
			break LOOP
		}
	}

}
