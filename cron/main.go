package main

import (
	"fmt"
	"sync"
	"time"

	"./cron"
)

/*
           *********
  grab the v2 branch to be able to Remove an entry from being run
  see: https://github.com/robfig/cron
*/

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()
	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.AddFunc("@every 10s", func() { fmt.Println(time.Now(), "Every 10 seconds") })
	c.AddFunc("@every 1m", func() { fmt.Println(time.Now(), "Every minute") })
	c.Start()
	defer c.Stop()

	fmt.Printf("%v -- %d entries=%T=%#v\n", time.Now(), len(c.Entries()), c.Entries(), c.Entries())
	fmt.Println(c.Entries())
	c.AddFunc("@every 2s", func() { fmt.Println(time.Now(), "Every 2 seconds ... after cron start!") })
	fmt.Printf("%v -- %d entries=%T=%+v\n", time.Now(), len(c.Entries()), c.Entries(), c.Entries()[0].Next)

	wg.Wait()
}
