package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
)

func main() {
    // response, _, err := http.Get("http://cellipede.com:4235/")
    response, err := http.Get("http://cellipede.com:4235/")
    if err != nil {
        fmt.Printf("Error: %s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        fmt.Printf("err: %s", err)
        if err != nil {
            fmt.Printf("Error: %s", err)
            os.Exit(1)
        }
        fmt.Printf("%s\n", string(contents))
    }
    fmt.Printf("\nerr: %s\n", err)
}
