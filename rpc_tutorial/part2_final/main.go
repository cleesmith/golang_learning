package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	for {
		// block until there's a connection to accept
		if conn, err := ln.Accept(); err == nil {
			go handleConnection(conn)
		}
	}
}

// handleConnection is spawned once per connection from a client,
// and exits when the client is done sending requests
func handleConnection(conn net.Conn) {
	// always close the conn no matter how this function exits:
	defer conn.Close()
	reader := bufio.NewReader(conn)
	// loop to handle multiple requests from user conn,
	// since we are now running concurrently:
	for {
		req, err := http.ReadRequest(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read request: %s", err)
			}
			return
		}
		// connect to a backend and send the request along
		if be, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
			be_reader := bufio.NewReader(be)
			if err := req.Write(be); err == nil {
				if resp, err := http.ReadResponse(be_reader, req); err == nil {
					FixHttp10Response(resp, req)
					if err := resp.Write(conn); err == nil {
						log.Printf("proxied %s: got %d", req.URL.Path, resp.StatusCode)
					}
					if resp.Close {
						return
					}
				}
			}
		}
	}
}
