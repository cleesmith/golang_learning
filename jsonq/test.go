package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/jsonq"
)

func main() {

	jsonstring := `
    {
      "foo": 1,
      "bar": 2,
      "test": "Hello, world!",
      "baz": 123.1,
      "array": [
        {"foo": 1},
        {"bar": 2},
        {"baz": 3}
      ],
      "subobj": {
        "foo": 1,
        "test2": "Hello, world2!",
        "subarray": [1,2,3],
        "subsubobj": {
          "bar": 2,
          "baz": 3,
          "array": ["hello", "world"]
        }
      },
      "bool": true
    }
  `
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(jsonstring))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	// fmt.Println(jq)
	// avalue, err := jq.Int("bar")
  // avalue, err := jq.String("tes") // an error
  // avalue, err := jq.String("test")
  avalue, err := jq.String("subobj", "test2")
	if err != nil {
		fmt.Printf("err=%s\n", err)
	} else {
		fmt.Printf("test=%v\n", avalue)
	}
  obj, err := jq.Object("subobj")
  fmt.Println(obj["test2"])
}
