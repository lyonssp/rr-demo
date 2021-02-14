package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// generate 10 random numbers in parallel and print them to stdout
func main() {
	rand.Seed(time.Now().UnixNano())

	n := 10
	ch := produce(n)

	for x := range ch {
		fmt.Println(fmt.Sprintf("%d", x))
	}
}

func produce(n int) <-chan int {
	ch := make(chan int, n)
	go func() {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				x := rand.Intn(100)
				ch <- x
				wg.Done()
			}()
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}
