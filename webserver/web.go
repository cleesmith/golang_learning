package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const host string = "localhost"
	const port string = "80"
	http.Handle("/", http.FileServer(http.Dir(".")))
	// if there's a file named "index.hmtl" in this directory (".")
	// then the above code will show it, otherwise it will list the directory files
	fmt.Printf("Listening on %s ...\n", host+":"+port)
	log.Output(http.ListenAndServe(host+":"+port, nil))
}
