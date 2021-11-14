package main

import (
  "fmt"
  "sync"
)

func gen(number, quantity int) <-chan int {
  ch := make(chan int)
  go func() {
    for i := 0; i < quantity; i++ {
      ch <- number
    }
    close(ch)
  }()
  return ch
}

func merge(chans ...<-chan int) <-chan int {
  result := make(chan int)
  wg := new(sync.WaitGroup)
  for _, ch := range chans {
    wg.Add(1)
    localChan := ch
    go func() {
      defer wg.Done()
      for num := range localChan {
        result <- num
      }
    }()
  }
  go func() {
    wg.Wait()
    close(result)
  }()
  return result
}

func main() {
  result := merge(gen(1, 10), gen(20, 2))
  for num := range result {
    fmt.Print(num, " ")
  }
}
