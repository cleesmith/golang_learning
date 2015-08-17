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
// via any youtube video page:
// 'subscribers', "//*[@id="watch7-subscription-container"]/span/span[@class='yt-subscription-button-subscriber-count-branded-horizontal ']/text()"
// via youtube user's about page:
// 'subscribers', "//*[@id='browse-items-primary']/li/div/div/div/span[@class='about-stat']/b/text()"
// 'subscribers', r"subscriber-count.*?>(?P<scrape>[0-9,]+?)<" ... old regex way
// 'ga',          r"(?:\'|\")(?P<scrape>UA-.*?)(?:\'|\")"

func main() {
	// resp, err := http.Get("http://amazon.com/")
	// resp, err := http.Get("http://cellipede.com:4235/")
	// resp, err := http.Get("http://cleesmith.github.io/")
	// resp, err := http.Get("https://github.com/cleesmith")
	// resp, err := http.Get("https://www.youtube.com/user/cleesmith2006/about")
	resp, err := http.Get("https://www.youtube.com/watch?v=Eacoqt4BtMc")
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
	// xp := xpath.Compile("//meta[translate(@name, 'ABCDEFGHJIKLMNOPQRSTUVWXYZ', 'abcdefghjiklmnopqrstuvwxyz')='description']/@content")

	// youtube user's about page, get both subscribers and total views:
	// xp := xpath.Compile("//*[@id='browse-items-primary']/li/div/div/div/span[@class='about-stat']/b/text()")

	// get subscribers via any youtube video page:
	xp := xpath.Compile("//*[@id='watch7-subscription-container']/span/span[@class='yt-subscription-button-subscriber-count-branded-horizontal ']/text()")

	nodes, err := doc.Root().Search(xp)
	if err != nil {
		fmt.Printf("ERROR: doc.Root().Search(xp): %v\n", err)
		panic(err)
		return
	}
	fmt.Printf("nodes=%T=%v=%v\n", nodes, len(nodes), nodes)
	if len(nodes) > 0 {
		subscribers := nodes[0].String()
		fmt.Printf("nodes[0]=%T=%v \nsubscribers=%T=%s\n", nodes[0], nodes[0].String(), subscribers, subscribers)
		// totalViews := nodes[1].String()
		// fmt.Printf("nodes[1]=%T=%v \ntotalViews=%T=%s\n", nodes[1], nodes[1].String(), totalViews, totalViews)
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
}
