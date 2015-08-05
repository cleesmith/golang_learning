package main

import "fmt"

type Counter int

func (c *Counter) AddOne() {
  *c++
}

func main() {
  var hits Counter
  hits.AddOne()
  fmt.Println(hits)
}
