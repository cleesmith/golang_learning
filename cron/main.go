package main

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/robfig/cron.v2"
)

/*
	use import "gopkg.in/robfig/cron.v2" to grab
  the v2 branch ... to Remove a job entry from being run
  see: https://github.com/robfig/cron
*/

type testJob struct {
	name string
}

func (t testJob) Run() {
	fmt.Printf("\tRun: %v ... job=%#v\n", time.Now(), t)
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()

	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.AddFunc("@every 1m", func() { fmt.Println(time.Now(), "Every minute") })
	c.AddFunc("@every 5s", func() { fmt.Println(time.Now(), "Every 5 seconds") })

	id_3s, err := c.AddJob("@every 3s", testJob{"every3s"})
	fmt.Printf("\nevery3s: job id=%v err=%v\n", id_3s, err)

	inspect(c.Entries())
	c.Start()
	defer c.Stop()

	fmt.Printf("\n%v -- %d entries=%#v\n", time.Now(), len(c.Entries()), c.Entries())
	fmt.Println(c.Entries())

	// can add/remove jobs/funcs after cron has started ...

	id_2s, err := c.AddFunc("@every 2s", func() { fmt.Println(time.Now(), "Every 2 seconds ... after cron start!") })
	fmt.Printf("\n*** job id=%v err=%v\n", id_2s, err)
	fmt.Printf("\n%v -- %d entries=%T=%+v\n", time.Now(), len(c.Entries()), c.Entries(), c.Entries()[0].Next)

	time.Sleep(10 * time.Second)
	c.Remove(id_2s)
	c.Remove(id_3s)

	wg.Wait()
}
