package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
)

type Movie struct {
  Title  string `json:"title"`
  Rating string `json:"rating"`
  Year   string `json:"year"`
}

// try it:
// go run json_api.go
// curl -X GET http://127.0.0.1:8080/movies -v
func main() {
  router := mux.NewRouter()
  router.HandleFunc("/movies", handleMovies).Methods("GET")
  http.ListenAndServe(":8080", router)
}

func handleMovies(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "application/json")
  var movies = map[string]*Movie{
    "tt0076759": &Movie{Title: "Star Wars: A New Hope", Rating: "8.7", Year: "1977"},
    "tt0082971": &Movie{Title: "Indiana Jones: Raiders of the Lost Ark", Rating: "8.6", Year: "1981"},
  }
  outgoingJSON, error := json.Marshal(movies)
  if error != nil {
    log.Println(error.Error())
    http.Error(res, error.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Fprint(res, string(outgoingJSON))
}
