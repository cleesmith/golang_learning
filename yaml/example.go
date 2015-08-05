package main

import (
  "fmt"
  "log"

  "gopkg.in/yaml.v2"
)

// var data = `
// a: Easy!
// b:
//   c: 2
//   d: [3, 4]
// `

// var data = `
// ---
// big_seo_client:
//   # this sample does everything, but you may do less of the following
//   a: Easy!
//   b:
//     c: 2
//     d:
//     - 3
//     - 4
// `

var data = `
---
big_seo_client:
  # this sample does everything, but you may do less of the following
  # note: be careful with indentation as yaml is picky about that, so
  #       I use 2 spaces and avoid the tab key
  project_description: |
    This is our big SEO client so take care of them.
    Ensure they know that SEO is not dead, so long as we're on the job.
    By "job" we mean wasting time and getting paid for doing very little.
    We should rename our company to "much ado about nothing".
  urls:
  - http://stackoverflow.com
  - http://www.amazon.com
  youtube_channels:
  - cleesmith2006
  - RomanAtwoodVlogs
  youtube_videos:
  - JFRNgYeA5WA
  - _LuWUZ1NbHQ
  twitter_accounts:
  - cleesmith
  - romanatwood
  serp_queries:
  - cellipede video
  - roman atwood
`

type T struct {
  A string
  B struct{C int; D []int ",flow"}
}

func main() {
  // t := T{}
  // err := yaml.Unmarshal([]byte(data), &t)
  // if err != nil {
  //   log.Fatalf("error: %v", err)
  // }
  // fmt.Printf("--- t:\n%v\n\n", t)
  // d, err := yaml.Marshal(&t)
  // if err != nil {
  //   log.Fatalf("error: %v", err)
  // }
  // fmt.Printf("--- t dump:\n%s\n\n", string(d))

  m := make(map[interface{}]interface{})

  err := yaml.Unmarshal([]byte(data), &m)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
  fmt.Printf("--- m:\n%v\n\n", m)
  fmt.Printf("type of m=%T\n", m)
  for k, v := range m {
    fmt.Printf("\nkey=%T=%v\n", k, k)
    fmt.Printf("\nvalue=%T=%v\n", v, v)
  }

  d, err := yaml.Marshal(&m)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
  fmt.Printf("\n\n--- m dump:\n%s\n\n", string(d))
}
