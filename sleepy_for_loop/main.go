package main

import (
	"fmt"
	"time"
	"unsafe"
)

func main() {
	i64 := int64(9223372036854775807)
	s := "12345678901234567890123456789012345678901234567890123456789012345678901234567890"
	empty_struct := struct{}{}

	// this "i" is not the same one as in the "for" loop below
	var i int = 9223372036854775807
	// in other words: "the scope of the variable"
	// According to the language specification “Go is lexically scoped using blocks”.
	// Basically this means that the variable exists within the nearest curly braces { } (a block)
	// including any nested curly braces (blocks), but not outside of them.

	for i := 0; i < 3; time.Sleep(3 * time.Second) {
		// see: https://medium.com/@felixge/the-sleepy-for-loop-in-go-4e6fee88c5ad#.x891gbs0r
		// a "sleepy for loop" with 2 retries
		fmt.Printf("trying: %v at: %v\n", i, time.Now())
		i++
		if i == 2 {
			fmt.Printf("done: i=%v at: %v\n", i, time.Now())
			break
		}
	}
	fmt.Printf("i=%v %v\n", i, time.Now())
	fmt.Printf("i=%v i64=%v s=%v\n", unsafe.Sizeof(i), unsafe.Sizeof(i64), unsafe.Sizeof(s))
	fmt.Printf("empty_struct=%T=%v\n", empty_struct, unsafe.Sizeof(empty_struct))
}
