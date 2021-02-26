package KlaudiasGoTrace

import (
	"math"
	"time"
)

func mySomething2(ch chan int, parentId uint64) {
	StartGoroutine(parentId)

	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 10000000; i++ {
		j := math.Sqrt(float64(i * i))
		j *= j * float64(i)
		_ = j
	}

	SendToChannel(42, ch)
	ch <- 42

	time.Sleep(10 * time.Millisecond)

	StopGoroutine()
}

func my2() {
	StartTrace()

	time.Sleep(10 * time.Millisecond)
	ch := make(chan int)

	go mySomething2(ch, GetGID())

	go func() {

	}()

	ReceiveFromChannel(<-ch, ch)

	time.Sleep(100 * time.Millisecond)

	EndTrace()
}
