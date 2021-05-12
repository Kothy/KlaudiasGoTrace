package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"fmt"
)

const MAX = 30

func Generate(ch chan<- int, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	for i := 2; ; i++ {
		KlaudiasGoTrace.SendToChannel(i, ch)

		ch <- i
	}
	KlaudiasGoTrace.StopGoroutine()

}

func Filter(in <-chan int, out chan<- int, prime int, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	for {
		i := <-in
		KlaudiasGoTrace.ReceiveFromChannel(i, in)

		if i%prime != 0 {
			KlaudiasGoTrace.SendToChannel(i, out)

			out <- i
		}
	}
	KlaudiasGoTrace.StopGoroutine()

}

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	ch := make(chan int)
	go Generate(ch, KlaudiasGoTrace.GetGID())
	for i := 0; i < MAX; i++ {
		prime := <-ch
		KlaudiasGoTrace.ReceiveFromChannel(prime, ch)

		fmt.Println(prime)
		ch1 := make(chan int)
		go Filter(ch, ch1, prime, KlaudiasGoTrace.GetGID())
		ch = ch1
	}
	KlaudiasGoTrace.EndTrace()

}
