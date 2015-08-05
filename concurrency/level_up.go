package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Printf("GOMAXPROCS=%v\n", runtime.GOMAXPROCS(-1))
	fmt.Printf("NumCPU=%v\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("GOMAXPROCS=%v\n", runtime.GOMAXPROCS(-1))
	max := 10
	intChan := make(chan int)

	for i := 1; i <= max; i++ {
		go func(j, max int) {
			fmt.Printf("go: %d\n", j)
			if j > 7 {
				fmt.Println("sleep 2 secs")
				time.Sleep(2 * time.Second) // work more
			} else {
				time.Sleep(1 * time.Millisecond)
			}
			intChan <- j
			if j >= max {
				close(intChan)
			}
		}(i, max)
	}

	for j := range intChan {
		fmt.Printf("main: %d\n", j)
	}
	fmt.Printf("GOMAXPROCS=%v\n", runtime.GOMAXPROCS(-1))
	fmt.Println("fin")
}
