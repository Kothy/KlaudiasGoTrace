package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"time"
)

func goHello(ch chan int, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)
	time.Sleep(10 * time.Millisecond)
	val := 42
	KlaudiasGoTrace.SendToChannel(val, ch)

	ch <- val
	time.Sleep(10 * time.Millisecond)
	KlaudiasGoTrace.StopGoroutine()

}

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	ch := make(chan int)

	go goHello(ch, KlaudiasGoTrace.GetGID())
	KlaudiasGoTrace.ReceiveFromChannel(<-ch, ch)
	KlaudiasGoTrace.EndTrace()

}
