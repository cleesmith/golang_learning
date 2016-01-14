package main

import (
	"fmt"
	"regexp"
)

type MatchResultResponse struct {
	Matches [][]string
	// GroupsName []string
}

func main() {
	matched, err := regexp.MatchString("^sea", "#seafood")
	fmt.Println(matched, err)
	matched, err = regexp.MatchString("bar.*", "seafood")
	fmt.Println(matched, err)
	// matched, err = regexp.MatchString("(ab", "seafood")
	// fmt.Printf("matched=%v \t err=%v\n", matched, err)

	r, err := regexp.Compile("^alert|^drop")
	matched = r.MatchString("# drop now!")
	fmt.Printf("matched=%v \t err=%v\n", matched, err)

	aline := "alert (msg:\"SQL sp_start_job - program execution\"; sd:123;)"

	r, err = regexp.Compile(`sid\s*:\s*(\d+);`)
	matched = r.MatchString(aline)
	fmt.Printf("matched=%v \t err=%v\n", matched, err)

	//Evaluate regex and get results
	matches := r.FindAllStringSubmatch(aline, -1)
	fmt.Printf("matches=%T=%#v\n", matches, matches)
	for _, m := range matches {
		for i2, m2 := range m {
			fmt.Printf("i2=%T=%#v \t m2=%T=%#v\n", i2, i2, m2, m2)
		}
	}
	result := &MatchResultResponse{}
	if len(matches) > 0 {
		result.Matches = matches
		// result.GroupsName = r.SubexpNames()[1:]
	}
	fmt.Printf("result=%T:\n%#v\n", result, result)

	// ************
	// much simpler
	// ************

	// aResult := r.FindAllStringSubmatch(aline, -1)
	aResult := r.FindStringSubmatch(aline)
	fmt.Printf("\naResult=%T=%#v\n", aResult, aResult)
	if aResult == nil {
		fmt.Println("NO sid found!")
	}
	for k, v := range aResult {
		fmt.Printf("%d. %s\n", k, v)
	}

	msgRx, err := regexp.Compile(`msg\s*:\s*\"(.*?)\";`)
	aResult = msgRx.FindStringSubmatch(aline)
	fmt.Printf("\naResult=%T=%#v\n", aResult, aResult)
	if aResult == nil {
		fmt.Println("NO msg found!")
	}
	for k, v := range aResult {
		fmt.Printf("%d. %s\n", k, v)
	}
}
