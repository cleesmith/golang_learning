package main

import (
	// "encoding/xml"
	// "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	// "net/url"
)

func main() {
	// aUrl := "http://www.amazon.com/"
	aUrl := "http://cleesmith.github.io/"
	resp, err := http.Get(aUrl)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}
	html_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}
	fmt.Printf("resp:\n%v\n", resp)
	fmt.Printf("html:\n%v\n", string(html_body))

	// xml_body := map[string]interface{}{}
	// err = xml.Unmarshal([]byte(html_body), &xml_body)
	// if err != nil {
	// 	fmt.Printf("ERROR: %v\n", err)
	// 	panic(err)
	// 	return
	// }
	// fmt.Printf("xml_body: %#v\n", xml_body)
}
