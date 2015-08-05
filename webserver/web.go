package main

import (
	"fmt"
	"net/http"
)

func main() {
  const host string = "localhost"
	const port string = "3000"
	http.Handle("/", http.FileServer(http.Dir(".")))
	// if there's a file named "index.hmtl" in this directory (".")
	// then the above code will show it, otherwise it will list the directory files
  fmt.Printf("Listening on %s ...\n", host + ":" + port)
	http.ListenAndServe(host + ":" + port, nil)
}
