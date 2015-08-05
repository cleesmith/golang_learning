package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/jmoiron/jsonq"
)

type Twitter struct {
	Count *big.Int `json:"count"`
	Url   string   `json:"url"`
}

func main() {
	value := url.Values{}
	// value.Add("url", "") // error msg from twitter api
	// value.Add("url", "http://www.spudamazon.com/") // 0 count
	value.Add("url", "http://www.amazon.com/")
	qstr := value.Encode()
	fmt.Println(qstr)
	url := "http://urls.api.twitter.com/1/urls/count.json?" + qstr
	fmt.Println(url)
	ch := make(chan string)
	go func() {
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ch <- string(body)
		close(ch)
	}()
	body := <-ch // blocks waiting for a message on this channel
	fmt.Printf("body=%s\n", body)
	var result Twitter
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("json Unmarshal error:", err)
	}
	fmt.Printf("result=%s\n", result)
	fmt.Printf("result.Count=%v\n", result.Count)

	// use jsonq:
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(body))
	err = dec.Decode(&data)
	if err != nil {
		fmt.Println("json Decode error:", err)
	}
	jq := jsonq.NewQuery(data)
	fmt.Println(jq)

	// // We need to provide a variable where the JSON
	// // package can put the decoded data. This
	// // `map[string]interface{}` will hold a map of strings
	// // to arbitrary data types. Unfortunately, this means all
	// // numbers will be float64, so it's better to use a known
	// // struct like "type Twitter struct" when the return json
	// // is of a known layout.
	// var dat map[string]interface{}
	// // here's the actual decoding, and a check for associated errors:
	// err = json.Unmarshal([]byte(body), &dat)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// // if err := json.Unmarshal([]byte(body), &dat); err != nil {
	// //  panic(err)
	// // }
	// fmt.Println(dat)
	// // In order to use the values in the decoded map,
	// // we'll need to cast them to their appropriate type.
	// // num := dat["count"].(float64)
	// num := dat["count"]
	// // fmt.Println(num.(int))
	// fmt.Println(num)
}
