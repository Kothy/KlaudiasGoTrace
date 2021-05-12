package main

import (
	"fmt"
	"time"
)

var number = 20

func main() {
	start := time.Now()
	prev := make(chan int)
	first := prev
	for i := 0; i < number; i++ {
		next := make(chan int)
		go func(from, to chan int) {
			for {
				val := <-from
				val += 1
				to <- val
			}
		}(prev, next)
		prev = next
	}
	elapsed := time.Since(start)
	first <- 0
	<-prev
	fmt.Println(elapsed)
}
