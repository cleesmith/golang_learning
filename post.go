package main

import (
    "bytes"
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {
    url := "http://cellipede.com:4235/"
    fmt.Println("URL:", url)

    var query = []byte(`your query`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "text/plain")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("\nresponse Status:\n", resp.Status)
    fmt.Println("\nresponse Headers:\n", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("\nresponse Body:\n", string(body))
}
