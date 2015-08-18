package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Backend struct {
	net.Conn // embedded type, i.e. a wrapper around net.Conn
	Reader   *bufio.Reader
	Writer   *bufio.Writer
}

var backendQueue chan *Backend
var requestBytes map[string]int64
var requestLock sync.Mutex

func init() {
	requestBytes = make(map[string]int64)
	backendQueue = make(chan *Backend, 10)
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	for {
		if conn, err := ln.Accept(); err == nil {
			go handleConnection(conn)
		}
	}
}

// handleConnection is spawned once per connection from a client,
// and exits when the client is done sending requests
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := http.ReadRequest(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read request: %s", err)
			}
			return
		}

		be, err := getBackend()
		if err != nil {
			return
		}

		if err := req.Write(be.Writer); err == nil {
			be.Writer.Flush()
			if resp, err := http.ReadResponse(be.Reader, req); err == nil {
				bytes := updateStats(req, resp)
				resp.Header.Set("X-Bytes", strconv.FormatInt(bytes, 10))

				FixHttp10Response(resp, req)
				if err := resp.Write(conn); err == nil {
					log.Printf("proxied %s: got %d", req.URL.Path, resp.StatusCode)
				}
				if resp.Close {
					return
				}
			}
		}

		go queueBackend(be)
	}
}

// updateStats takes a request and response and collects some statistics about them
func updateStats(req *http.Request, resp *http.Response) int64 {
	requestLock.Lock()
	defer requestLock.Unlock()
	bytes := requestBytes[req.URL.Path] + resp.ContentLength
	requestBytes[req.URL.Path] = bytes
	return bytes
}

// either gets a backend from the queue or makes a new backend
func getBackend() (*Backend, error) {
	select {
	case be := <-backendQueue:
		return be, nil
	case <-time.After(100 * time.Millisecond):
		be, err := net.Dial("tcp", "127.0.0.1:8081")
		if err != nil {
			return nil, err
		}
		return &Backend{
			Conn:   be,
			Reader: bufio.NewReader(be),
			Writer: bufio.NewWriter(be),
		}, nil
	}
}

// takes a backend and re-enqueues it
func queueBackend(be *Backend) {
	select {
	case backendQueue <- be:
		// nothing to do here, as the backend re-enqueued safely
	case <-time.After(1 * time.Second):
		// If this fires, it means the queue of backends was full.
		// We don't want to wait forever, as this period of time
		// blocks us handling the next request a user might send us.
		be.Close()
	}
}
