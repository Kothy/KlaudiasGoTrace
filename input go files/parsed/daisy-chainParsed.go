package main

import "KlaudiasGoTrace/KlaudiasGoTrace"

func f(left, right chan int, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)
	val := <-right
	KlaudiasGoTrace.ReceiveFromChannel(val, right)
	val += 1
	KlaudiasGoTrace.SendToChannel(val, left)

	left <- val
	KlaudiasGoTrace.StopGoroutine()

}

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	const n = 100
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right, KlaudiasGoTrace.GetGID())
		left = right
	}
	go func(c chan int, parentId uint64) {
		KlaudiasGoTrace.StartGoroutine(parentId)
		KlaudiasGoTrace.SendToChannel(1, c)
		c <- 1
		KlaudiasGoTrace.StopGoroutine()

	}(right, KlaudiasGoTrace.GetGID())
	KlaudiasGoTrace.ReceiveFromChannel(<-leftmost,
		leftmost)
	KlaudiasGoTrace.EndTrace()

}
