package main

import "fmt"

// both send/receive channel: "c chan int"
// receive only channel: "c <-chan int"
// send only channel: "c chan<- int"
func sum(a []int, c chan<- int) {
	sum := 0
	for _, v := range a {
		sum += v
	}
	c <- sum // send sum to c

	// if channel is both send/receive:
	// d := <-c
	// fmt.Println(d)
	// but this causes:
	//  fatal error: all goroutines are asleep - deadlock!
}

func main() {
	a := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)

	go sum(a[:len(a)/2], c)
	fmt.Println(a[:len(a)/2])

	go sum(a[len(a)/2:], c)
	fmt.Println(a[len(a)/2:])

	// note: 2 goroutines so 2 channel receives:
	x, y := <-c, <-c // receive from c

	fmt.Println(x, y, x+y)
}
