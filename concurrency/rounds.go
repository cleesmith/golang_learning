package main

func main() {
	for i := 0; i < 5; i++ {
		println("i=", i)
		var (
			first  = make(chan struct{}, 1)
			second = make(chan struct{}, 1)
		)
		first <- struct{}{}
		second <- struct{}{}
		select {
		case <-first:
			println("round 1: first")
		case <-second:
			println("round 1: second")
		}
		select {
		case <-first:
			println("round 2: first")
		case <-second:
			println("round 2: second")
		}
	}
}
