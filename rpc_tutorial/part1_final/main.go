package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

func main() {
	// listen for connections
	if ln, err := net.Listen("tcp", ":8080"); err == nil {
		for {
			// accept connections
			if conn, err := ln.Accept(); err == nil {
				reader := bufio.NewReader(conn)
				// read requests from the client
				if req, err := http.ReadRequest(reader); err == nil {
					// connect to the backend web server
					if be, err := net.Dial("tcp", "127.0.0.1:8081"); err == nil {
						be_reader := bufio.NewReader(be)
						// send the request to the backend
						if err := req.Write(be); err == nil {
							// read the response from the backend
							if resp, err := http.ReadResponse(be_reader, req); err == nil {
								// send the response to the client, making sure to close it
								resp.Close = true
								if err := resp.Write(conn); err == nil {
									log.Printf("proxied %s: got %d", req.URL.Path, resp.StatusCode)
								}
								conn.Close()
								// loop back to accepting the next connection
							}
						}
					}
				}
			}
		}
	}
}
