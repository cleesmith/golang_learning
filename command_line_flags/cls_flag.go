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

	// see: https://gist.github.com/rcrowley/5474430
	//      for handling SIGINT and SIGTERM
	// handle Ctrl+C interrupt:
	exitChan := make(chan os.Signal, 1)
	// the following seems to work for:
	//   The only signal values guaranteed to be present on all systems are
	//   Interrupt (send the process an interrupt) and
	//   Kill (force the process to exit).
	//   quote from: https://golang.org/pkg/os/#Signal
	signal.Notify(exitChan, os.Kill, os.Interrupt)
	defer signal.Stop(exitChan)

	for {
		select {
		// case <-time.After(1 * time.Second):
		// 	panicIf(blink(false))
		// 	// panicIf(blink(true))
		case <-exitChan:
			fmt.Println("\nctrl+c pressed")
			return
		default:
			// do
			fmt.Println("default")
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
	// fmt.Printf("\npanicIf: err=%v\n", err)
	if err != nil {
		panic(err)
	}
}
