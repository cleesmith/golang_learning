package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOMAXPROCS(-1)) //prints: 1

	fmt.Println(runtime.NumCPU()) //prints: 1 (on play.golang.org)

	// this will display what the "ulimit -n" might be:
	runtime.GOMAXPROCS(100000000)
	fmt.Println(runtime.GOMAXPROCS(-1)) //prints: 256

	runtime.GOMAXPROCS(1)               // set to use just 1 cpu
	fmt.Println(runtime.GOMAXPROCS(-1)) //prints: 1
}
