package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
)

type Tweet struct {
	// use a known type of "big.Int" to covert json's
	// "number"type (a float64 by default), but we
	// don't want a floating point number coz it's a count
	Count *big.Int `json:"count"`
	Url   string   `json:"url"`
}

func main() {
	// tweet, err := getTweets("")
	// tweet, err := getTweets("http://www.spudamazon.com/")
	tweet, err := getTweets("http://www.amazon.com/")
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}
	fmt.Printf("tweet=%s\n", tweet)
	fmt.Printf("tweet.Count=%v\n", tweet.Count)
}

func getTweets(aUrl string) (Tweet, error) {
	var tweet Tweet
	var err error
	value := url.Values{}
	value.Add("url", aUrl)
	qstr := value.Encode()
	// this should be a config constant somewhere, in case it changes:
	twitterUrl := "http://urls.api.twitter.com/1/urls/count.json?" + qstr
	resp, _ := http.Get(twitterUrl)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tweet, errors.New("getTweets: " + err.Error())
	}
	if err = json.Unmarshal(body, &tweet); err != nil {
		return tweet, errors.New("getTweets: " + err.Error())
	}
	return tweet, err
}
