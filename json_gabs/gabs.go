package main

import (
	"fmt"
	"strconv"

	"github.com/jeffail/gabs" // for dealing with dynamic or unknown JSON structures
)

func main() {
	// typical twitter json response:
	// jsonParsed, err := gabs.ParseJSON([]byte(`{"count":84880915,"url":"http:\/\/www.amazon.com\/"}`))
	// fmt.Printf("jsonParsed=%v err=%v\n", jsonParsed, err)

	// twitter json response for a bad request:
	jsonParsed, err := gabs.ParseJSON([]byte(`{"request":"\/1\/urls\/count.json?url=","error":"Missing or invalid url parameter."}`))
	fmt.Printf("jsonParsed=%v err=%v\n", jsonParsed, err)

	var url string
	var value float64
	var ok bool

	url, ok = jsonParsed.Path("url").Data().(string)
	fmt.Printf("url=%T=%v ok=%v\n", url, url, ok)

	value, ok = jsonParsed.Path("count").Data().(float64)
	fmt.Printf("value=%T=%v ok=%v\n", value, value, ok)
	if ok {
		fmt.Println("it's OK, coz we just put the string in a cell")
		count := strconv.FormatFloat(value, 'f', 0, 64)
		fmt.Printf("count=%T=%v\n", count, count)
	}
}
