package main

import (
	// "errors" // to test panic
	// "fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var Info *log.Logger

var chttp = http.NewServeMux()

func main() {
	Info = log.New(os.Stdout,
		"",
		// log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds|log.Llongfile)
		log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":80", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	// err = errors.New("spud") // test panic
	if err != nil {
		Info.Println("*** Error: panicking:")
		panic(err.Error())
	}
	Info.Printf("%s \"%s\" for %s : %s\n", r.Method, r.URL.Path, host, port)
	chttp.ServeHTTP(w, r)
}
