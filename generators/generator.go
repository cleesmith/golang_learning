//~~~~~~~~~~~~~~~~~~~~~~~~~~~~//
// Generator that counts to n //
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~//
package main

import (
	"fmt"
)

func count(n int) chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func main() {
	for i := range count(100) {
		fmt.Println("Counted", i)
	}
}
