package main

import (
	"fmt"
	"golang.org/x/net/websocket" // go get golang.org/x/net/websocket
	"net/http"
)

func main() {
	http.Handle("/", websocket.Handler(handler))
	http.ListenAndServe("localhost:3000", nil)
}

func handler(c *websocket.Conn) {
	var s string
	fmt.Fscan(c, &s)
	fmt.Println("Received:", s)
	fmt.Fprint(c, "How do you do?")
}

// test using these in browser console:
// var sock = new WebSocket("ws://localhost:3000/");
// sock.onmessage = function(m) { console.log("Received:", m.data); }
// sock.send("Hello!\n")
