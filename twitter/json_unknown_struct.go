package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var data = `{
  "id": 12423434,
  "Name": "Fernando"
}`

func main() {
	d := json.NewDecoder(strings.NewReader(data))
	d.UseNumber()
	var x interface{}
	if err := d.Decode(&x); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decoded to %#v\n", x)
	result, err := json.Marshal(x)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("encoded to %s\n", result)
}
