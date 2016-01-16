package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var doStuffCh = make(chan struct{})
var aSignal os.Signal
var wg sync.WaitGroup

func main() {
	go doStuff()
	wg.Add(1) // so we can wait on doStuff to finish

	signalCh := make(chan os.Signal, 1)

	// quote from: https://golang.org/pkg/os/#Signal
	//   The only signal values guaranteed to be present on all systems are
	//   Interrupt (send the process an interrupt) and
	//   Kill (force the process to exit).
	// note: a "kill -9 pid" on OS X never reaches this program
	// note: a "kill -15 pid" is the same as "kill pid" (it's the default)
	//                    KILL -15 pid             ctrl+c  KILL -9 pid      KILL -1 pid
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
	// SIGHUP should cause us to re-read config files and start over
	defer func() {
		signal.Stop(signalCh)
		fmt.Printf("defer: wait for workers to finish: wg=%#v\n", wg)
		wg.Wait()
		fmt.Printf("defer: '%v' signal received! \tcalling signal.Stop!\n", aSignal)
	}()

	aSignal = <-signalCh
	fmt.Printf("\nmain: '%v' signal received\n", aSignal)

	fmt.Printf("\nmain: tell doStuff to quit by calling stopDoStuff()\n")
	stopDoStuff()
}

func doStuff() {
	defer wg.Done()
	for {
		select {
		case <-doStuffCh:
			fmt.Println("doStuff: 'case <-doStuffCh:' returning!")
			return
		default:
			fmt.Print(".")
		}
	}
}

func stopDoStuff() {
	fmt.Println("\nstopDoStuff: closing channel doStuffCh")
	close(doStuffCh)
	fmt.Println("\nstopDoStuff: closed channel doStuffCh!")
}
