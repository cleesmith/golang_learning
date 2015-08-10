package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	flag.Parse()

	// handle Ctrl+C interrupt:
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	for {
		select {
		case <-time.After(1 * time.Second):
			panicIf(blink(false))
			// panicIf(blink(true))
		case <-quit:
			fmt.Println("\nctrl+c pressed")
			return
		}
	}
	time.Sleep(10 * time.Second)
}

func blink(status bool) (err error) {
	fmt.Println("blink!")
	if status {
		err = errors.New("Error: unable to blink")
	}
	return err
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
