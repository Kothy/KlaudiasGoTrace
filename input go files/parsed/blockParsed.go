package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"math"
	"time"
)

func goBlocked(ch chan int, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 10000000; i++ {
		j := math.Sqrt(float64(i * i))
		j *= j * float64(i)
		_ = j
	}
	val := 42
	KlaudiasGoTrace.SendToChannel(val,
		ch)

	ch <- val

	time.Sleep(10 * time.Millisecond)
	KlaudiasGoTrace.StopGoroutine()

}

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	time.Sleep(10 * time.Millisecond)
	ch := make(chan int)
	go goBlocked(ch, KlaudiasGoTrace.GetGID())
	time.Sleep(100 * time.Millisecond)
	KlaudiasGoTrace.EndTrace()

}
