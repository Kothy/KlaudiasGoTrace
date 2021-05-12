package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"math/rand"
	"time"
)

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	var Ball int
	table := make(chan int)

	go player(table, KlaudiasGoTrace.GetGID())
	go player(table, KlaudiasGoTrace.GetGID())

	time.Sleep(20 * time.Millisecond)
	KlaudiasGoTrace.SendToChannel(Ball,
		table)

	table <- Ball
	time.Sleep(1 * time.Second)
	KlaudiasGoTrace.ReceiveFromChannel(<-table,
		table)
	KlaudiasGoTrace.EndTrace()

}

func player(table chan int, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	for {
		ball := <-table
		KlaudiasGoTrace.ReceiveFromChannel(ball,
			table)

		ball++
		n := (rand.Intn(100-50) + 50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		KlaudiasGoTrace.SendToChannel(ball,
			table)

		table <- ball
	}
	KlaudiasGoTrace.StopGoroutine()

}
