package main

import (
	// "errors" // to test panic
	// "fmt"
	"log"
	"net"
	"net/http"
	"os"
	// "strings"
	"time"
)

// typical ror log:
// 		Started GET "/" for 127.0.0.1 at 2014-08-03 22:31:32 -0400
// 		Processing by WelcomeController#index as HTML
// 		  Rendered welcome/index.html.erb (0.8ms)
// 		Completed 200 OK in 10ms (Views: 4.3ms | ActiveRecord: 0.0ms)

var Info *log.Logger

var chttp = http.NewServeMux()

func main() {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Println("err: ", err.Error())
	}
	log.Printf("utc=%v\ttime=%v\n", utc, time.Now().In(utc))

	log.Printf("flags=%v\n", log.Flags())
	// log.SetFlags(log.Ldate)
	// log.SetFlags(log.Ltime)
	// log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	log.SetFlags(log.LUTC)
	log.Printf("flags=%v\n", log.Flags())

	t, o := time.Now().Zone()
	log.Printf("Timezone: %v %v", t, o)
	Info = log.New(os.Stdout,
		"",
		log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds|log.Llongfile)
	// log.Ldate|log.Ltime|log.Lshortfile)
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", HomeHandler) // homepage
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
	// if strings.Contains(r.URL.Path, ".") {
	// 	chttp.ServeHTTP(w, r)
	// } else {
	// 	// fmt.Fprintf(w, "?\n")
	// 	chttp.ServeHTTP(w, r)
	// }
}
