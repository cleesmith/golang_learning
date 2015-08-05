package main

import (
  "fmt"
  "sync"
  "time"
)

func main() {
  var wg sync.WaitGroup
  done := make(chan struct{})
  wq := make(chan interface{})
  worker_count := 2

  for i := 0; i < worker_count; i++ {
    wg.Add(1)
    go doit(i, wq, done, &wg)
  }

  fmt.Printf("doing work\n")
  for i := 0; i < worker_count; i++ {
    time.Sleep(time.Millisecond * time.Duration(100))
    wq <- fmt.Sprintf("worker: %d", i)
  }

  fmt.Printf("closing 'done' channel\n")
  close(done)
  fmt.Printf("block/wait until all workers are done\n")
  wg.Wait()
  fmt.Println("all done!")
}

func doit(worker_id int, wq <-chan interface{}, done <-chan struct{}, wg *sync.WaitGroup) {
  fmt.Printf("[%v] is working\n", worker_id)
  defer wg.Done()
  max_time := time.Second * time.Duration(5)
  for {
    select {
    case m := <-wq:
      fmt.Printf("[%v] m => %v\n", worker_id, m)
    case <-done:
      fmt.Printf("[%v] is done\n", worker_id)
      return
    case <-time.After(max_time):
      fmt.Printf("timeout > %s seconds!\n", max_time)
    }
  }
}
