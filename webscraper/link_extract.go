package main

import (
	// "crypto/sha1"
	"errors"
	"fmt"
	// "io"
	// "net/http"
	"os"
	// "strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	// "github.com/Sirupsen/logrus"
	"github.com/cryptix/go/logging"
)

var (
	log    = logging.Logger("linkExt")
	hashes *os.File
)

// const url = `http://www.spiegel.de/international/germany/inside-the-nsa-s-war-on-internet-security-a-1010361.html`
const url = "http://cellipede.com:4235/"

func main() {
	var wg sync.WaitGroup

	doc, err := goquery.NewDocument(url)
	logging.CheckFatal(err)

	hashes, err = os.Create("hashes")
	logging.CheckFatal(err)
	defer hashes.Close()

	doc.Find("a").Each(
		func(i int, s *goquery.Selection) {
			link, found := s.Attr("href")
			title, _ := s.Attr("title")
			// if found && strings.HasSuffix(link, ".pdf") {
			if found {
				wg.Add(1)
				go fetchPDF(&wg, link, title)
			}
			return
		})
	wg.Wait()
	// log.Notice("Done")
	log.Println("Done")
}

func fetchPDF(wg *sync.WaitGroup, l, t string) (err error) {
	fmt.Printf("\nfetch: %v\n", l)
	// s := sha1.New()

	// fname := l[7:len(l)-4] + "-" + strings.TrimSpace(t) + ".pdf"
	// fname = strings.Replace(fname, "/", "-", -1)
	// fname = "pdfs/" + fname
	// log.Noticef("fetching: %s", fname)
	// log.Printf("fetching: %s", fname)

	defer func() {
		fmt.Println(">>> defer func <<<")
		if err != nil {
			fmt.Printf("\t calling myself: err=%v\n", err)
			fetchPDF(wg, l, t)
		} else {
			fmt.Println("\t no error, so we're done")
			// fmt.Fprintf(hashes, "%x %s\n", s.Sum(nil), fname)
			wg.Done()
		}
	}()

	err = errors.New("ERROR!!!")
	// resp, err := http.Get("https://www.spiegel.de/" + l)
	// if err != nil {
	// 	// log.Critical(err)
	// 	log.Println(err)
	// 	return
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	// log.Criticalf("http.Get %q", resp.Status)
	// 	log.Printf("http.Get %q", resp.Status)
	// 	return
	// }

	// f, err := os.Create(fname)
	// if err != nil {
	// 	// log.Critical(err)
	// 	log.Println(err)
	// 	return
	// }

	// multi := io.MultiWriter(s, f)
	// _, err = io.Copy(multi, resp.Body)
	// if err != nil {
	// 	// log.Critical(err)
	// 	log.Println(err)
	// 	os.Remove(fname)
	// 	return err
	// }
	// // log.Noticef("Saved: %s", fname)
	// log.Printf("Saved: %s", fname)

	return nil
	// return err
}
