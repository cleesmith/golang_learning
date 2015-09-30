package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	// byt := []byte(`{"StumbleUpon":0,"Reddit":"null","Facebook":{"commentsbox_count":14,"click_count":0,"total_count":0,"comment_count":0,"like_count":0,"share_count":0},"Delicious":0,"GooglePlusOne":843502,"Buzz":null,"Twitter":6597517,"Diggs":0,"Pinterest":119,"LinkedIn":99216}`)
	// {"Error":"Not a valid URL.","Type":"invalid_url","HTTP_Code":401}
	// byt := []byte(`{"Error":"Not a valid URL.","Type":"invalid_url","HTTP_Code":401}`)
	// {"Error":"Not a valid URL.", "Type": "invalid_url"}
	// {"Error": "Free API quota exceeded. Quota will reset tomorrow. Visit SharedCount to inquire about paid plans. http://sharedcount.com/quota", "Type": "quota_exceeded"}
	// byt := []byte(`{"Error": "Free API quota exceeded. Quota will reset tomorrow. Visit SharedCount to inquire about paid plans. http://sharedcount.com/quota", "Type": "quota_exceeded"}`)
	// {"Error": "Your Request uses jQuery\'s cachebusting function. Visit http://sharedcount.com/jsonp to get access to a cache-friendly JavaScript plugin you can use instead.", "Type":"no_jquery_cachebusting"}
	// note: \' is an "err=invalid character '\'' in string escape code" to "json.Unmarshal"
	byt := []byte(`{"Error": "Your Request uses jQuery's cachebusting function. Visit http://sharedcount.com/jsonp to get access to a cache-friendly JavaScript plugin you can use instead.", "Type":"no_jquery_cachebusting"}`)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		fmt.Printf("err=%v\n", err)
	}
	fmt.Printf("dat=%T=%+v\n", dat, dat)

	// tweets := strconv.FormatFloat(dat["Twitter"].(float64), 'f', 0, 64)
	// fmt.Printf("tweets=%T=%v\n", tweets, tweets)

	fb := dat["Facebook"]
	fmt.Printf("\nfb=%T=%+v\n", fb, fb)

	// printout from looping over "dat":
	// GooglePlusOne is float64=843502 count=843502
	// Diggs is float64=0  count=0
	// Pinterest is float64=119  count=119
	// LinkedIn is float64=99216 count=99216
	// StumbleUpon is float64=0  count=0
	// Facebook is a map[string]interface{}:
	//   like_count=0
	//   share_count=0
	//   commentsbox_count=14
	//   click_count=0
	//   total_count=0
	//   comment_count=0
	// Delicious is float64=0  count=0
	// Buzz is float64=0 count=0
	// Twitter is float64=6.597517e+06 count=6597517
	// Reddit is float64=0 count=0
	// k:  is the name of social count
	// vv: is the actual count as float64
	for k, v := range dat {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			count := strconv.FormatFloat(vv, 'f', 0, 64)
			fmt.Printf("%v is float64=%v\tcount=%v\n", k, vv, count)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case map[string]interface{}:
			fmt.Println(k, "is a map[string]interface{}:")
			for i, u := range vv {
				fmt.Printf("\t%v=%v\n", i, u)
			}
		default:
			fmt.Printf("%v is of a type I don't know how to handle: vv=%T=%v\n", k, vv, vv)
		}
	}
}
