package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// try it:
// go run main.go
// curl -X GET http://127.0.0.1:8080/hello/molly -v
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", index).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to /hello request")
	log.Println(r.UserAgent())
	vars := mux.Vars(r)
	name := vars["name"]
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}
