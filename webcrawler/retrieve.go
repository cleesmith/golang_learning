package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("http://cellipede.com:4235/")
	fmt.Println("http transport error is:", err)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("read error is:", err)
	fmt.Println(string(body))
}
