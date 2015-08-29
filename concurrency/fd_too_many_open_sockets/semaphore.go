// see: http://burke.libbey.me/conserving-file-descriptors-in-go/
// moderate File Descriptor consumption in two simple ways:
//  1. worker pool
//  2. semaphore
// the following is the semaphore strategy:
package main

const nTokens = 100

func main() {
  up := make(chan string)
  crawl(up)
}

func crawl(urlProducer chan string) {
  sem := make(chan bool, nTokens)

  for i := 0, i < nTokens; i++ {
    sem <- true
  }

  for url := range urlProducer {
    go func() {
      <- sem
      defer sem <- true

      fetch(<-urlProducer)
    }
  }()
}
