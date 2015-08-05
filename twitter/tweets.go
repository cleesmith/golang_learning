package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
)

func main() {
	f := 123.456
	fmt.Printf("f=%T\n", int(f))
	if false {
		tweets()
	}
}

type Twitter struct {
	// use a known type of "big.Int" to covert json's
	// type "number"(by default is float64), but we
	// don't want a floating point number
	Count *big.Int `json:"count"`
	Url   string   `json:"url"`
}

func tweets() {
	value := url.Values{}
	// value.Add("url", "") // error msg from twitter api
	// value.Add("url", "http://www.spudamazon.com/") // 0 count
	value.Add("url", "http://www.amazon.com/")
	qstr := value.Encode()
	fmt.Println(qstr)
	url := "http://urls.api.twitter.com/1/urls/count.json?" + qstr
	fmt.Println(url)

	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("body=%s\n", body)

	var result Twitter
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("json Unmarshal error:", err)
	}
	// if err := json.Unmarshal(body, &result); err != nil {
	//   return err
	// }
	fmt.Printf("result=%s\n", result)
	fmt.Printf("result.Count=%v\n", result.Count)

	// data := map[string]interface{}{}
	// dec := json.NewDecoder(strings.NewReader(body))
	// err = dec.Decode(&data)
	// if err != nil {
	// 	fmt.Println("json Decode error:", err)
	// }
	// jq := jsonq.NewQuery(data)
	// fmt.Println(jq)
}
