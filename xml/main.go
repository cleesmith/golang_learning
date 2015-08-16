package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"
)

// 'title',       "//title/text()"
// 'description', "//meta[translate(@name, 'ABCDEFGHJIKLMNOPQRSTUVWXYZ', 'abcdefghjiklmnopqrstuvwxyz')='description']/@content"
// 'canonical',   "/html/head/link[@rel = 'canonical']/@href"
// 'mobile',      "/html/head/link[@media = 'only screen and (max-width: 640px)']/@href"
// 'tweettotal',  "//span[.='Tweets']/following-sibling::span/text()"
// 'following',   "//span[.='Following']/following-sibling::span/text()"
// 'followers',   "//span[.='Followers']/following-sibling::span/text()"
// 'views',       "//div[@class='watch-view-count']/text()"
// 'thumbsup',    "//button[@id='watch-like']/span/text()"
// 'thumbsdown',  "//button[@id='watch-dislike']/span/text()"
// 'subscribers', r"subscriber-count.*?>(?P<scrape>[0-9,]+?)<"
// 'ga',          r"(?:\'|\")(?P<scrape>UA-.*?)(?:\'|\")"

func main() {
	// resp, err := http.Get("http://amazon.com/")
	// resp, err := http.Get("http://cellipede.com:4235/")
	resp, err := http.Get("http://cleesmith.github.io/")
	// resp, err := http.Get("https://github.com/cleesmith")
	if err != nil {
		fmt.Printf("ERROR: http.Get: %v\n", err)
		panic(err)
		return
	}
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: ioutil.ReadAll: %v\n", err)
		panic(err)
		return
	}
	// doc, err := gokogiri.ParseHtml([]byte(page))
	doc, err := gokogiri.ParseHtml(page)
	// important -- don't forget to free the resources:
	defer doc.Free()
	if err != nil {
		fmt.Printf("ERROR: gokogiri.ParseHtml: %v\n", err)
		panic(err)
		return
	}

	// perform operations on the parsed page -- consult the tests for examples
	fmt.Printf("page:\n%v\nresp=%v\n---------------\n", string(page), resp)
	// xp := xpath.Compile("/html/body/hr")
	// xp := xpath.Compile("//title/text()")
	xp := xpath.Compile("//meta[translate(@name, 'ABCDEFGHJIKLMNOPQRSTUVWXYZ', 'abcdefghjiklmnopqrstuvwxyz')='description']/@content")
	nodes, _ := doc.Root().Search(xp)
	fmt.Printf("nodes=%T=%v=%v\n", nodes, len(nodes), nodes)
	attrValue := nodes[0].String()
	fmt.Printf("nodes[0]=%T=%v \nattrValue=%T=%s\n", 
		nodes[0], nodes[0].String(), attrValue, attrValue)
	fmt.Println("\n- nodes matching search:")
	for n := range nodes {
		fmt.Printf("\tnodes[%v]=%T=%s nodes[%v].Name()=%v\n",
			n, nodes[n].InnerHtml(), nodes[n], n, nodes[n].Name())
		// subnodes, _ := nodes[n].Search("bar")
		// for s := range subnodes {
		// 	fmt.Println(subnodes[s].Name())
		// }
	}
	// for _, s := range ss {
	// 	ww, _ := s.Search(xpw)
	// 	for _, w := range ww {
	// 		fmt.Println(w.InnerHtml())
	// 	}
	// }
}
