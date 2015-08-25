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
	// jsonParsed, err := gabs.ParseJSON([]byte(`{"request":"\/1\/urls\/count.json?url=","error":"Missing or invalid url parameter."}`))
	// fmt.Printf("jsonParsed=%v err=%v\n", jsonParsed, err)

	// stumbleupon:
	jsonParsed, err := gabs.ParseJSON([]byte(`{ "result":
    { "url":"http:\/\/www.amazon.com\/","in_index":true,"publicid":"1ByuY3",
      "views":4958,
      "title":"Amazon.com: Online Shopping for Electronics, Apparel, Computers, Books, DVDs & more",
      "thumbnail":"http:\/\/cdn.stumble-upon.com\/mthumb\/788\/626788.jpg",
      "thumbnail_b":"http:\/\/cdn.stumble-upon.com\/bthumb\/788\/626788.jpg",
      "submit_link":"http:\/\/www.stumbleupon.com\/badge\/?url=http:\/\/www.amazon.com\/",
      "badge_link":"http:\/\/www.stumbleupon.com\/badge\/?url=http:\/\/www.amazon.com\/",
      "info_link":"http:\/\/www.stumbleupon.com\/url\/www.amazon.com\/"
    },"timestamp":1433685416,"success":true
  }`))
	// JSON.parse(response.body)['result']['views'] rescue response.body
	fmt.Printf("jsonParsed=%v err=%v\n", jsonParsed, err)

	// var url string
	var value float64
	var ok bool

	value, ok = jsonParsed.Path("result.views").Data().(float64)
	fmt.Printf("value=%T=%v ok=%v\n", value, value, ok)
	if ok {
		count := strconv.FormatFloat(value, 'f', 0, 64)
		fmt.Printf("count=%T=%v\n", count, count)
	}

	// url, ok = jsonParsed.Path("url").Data().(string)
	// fmt.Printf("url=%T=%v ok=%v\n", url, url, ok)

	// value, ok = jsonParsed.Path("count").Data().(float64)
	// fmt.Printf("value=%T=%v ok=%v\n", value, value, ok)
	// if ok {
	// 	count := strconv.FormatFloat(value, 'f', 0, 64)
	// 	fmt.Printf("count=%T=%v\n", count, count)
	// }
}
