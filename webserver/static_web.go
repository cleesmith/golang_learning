package main

import (
	// "errors" // to test panic
	// "fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

var chttp = http.NewServeMux()

func main() {
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", HomeHandler) // homepage
	http.ListenAndServe(":80", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	// err = errors.New("spud") // test panic
	if err != nil {
		log.Println("*** Error: panicking:")
		panic(err.Error())
	}
	log.Printf("host=%s port=%s\n", host, port)
	log.Printf("r.URL.Path=%v\n", r.URL.Path)
	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {
		// fmt.Fprintf(w, "?\n")
		chttp.ServeHTTP(w, r)
	}
}
