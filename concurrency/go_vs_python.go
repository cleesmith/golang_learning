// see python version:
// https://www.youtube.com/watch?v=MCs5OvhV9S4
// https://github.com/dabeaz/concurrencylive
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
)

func fib(n int64) int64 {
	if n <= 2 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func fibServer(addr string) {
	fmt.Printf("Listening on %v\n", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(fmt.Errorf("An error occured while listening to: %s -- %s", addr, err))
	}
	for {
		conn, err := ln.Accept()
		log.Println("connection", addr)
		if err != nil {
			log.Println(err)
			continue
		}
		go fibHandler(conn) // prefix by `go` to get concurrent.go
	}
}

func fibHandler(conn net.Conn) {
	buf := make([]byte, 100)
	var req int
	for {
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			conn.Close()
			break
		}
		reqStr := string(bytes.Trim(buf[0:n], "\n"))
		req, err = strconv.Atoi(reqStr)
		if err != nil {
			log.Println("The request must be a number", reqStr, err)
		}
		result := fmt.Sprintf("%v\n", fib(int64(req)))
		_, err = conn.Write([]byte(result))
		if err != nil {
			fmt.Println("Error while writing to the socket")
		}
	}
	log.Println("closed")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("You must provide an addr (127.0.0.1:25000)")
		return
	}
	fibServer(args[0])
}
