package KlaudiasGoTrace

import (
	"math"
	"time"
)

func myGoroutine(ch chan int) {
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 10000000; i++ {
		j := math.Sqrt(float64(i * i))
		j *= j * float64(i)
		_ = j
	}

	ch <- 42
	time.Sleep(10 * time.Millisecond)

}

func my() {
	time.Sleep(10 * time.Millisecond)
	ch := make(chan int)
	go myGoroutine(ch)
	<-ch
	time.Sleep(100 * time.Millisecond)

}
