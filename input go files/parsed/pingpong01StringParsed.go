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

	var Ball string
	Ball = "Ball"
	table := make(chan string)

	go playerString(table, KlaudiasGoTrace.GetGID())
	go playerString(table, KlaudiasGoTrace.GetGID())

	time.Sleep(20 * time.Millisecond)
	KlaudiasGoTrace.SendToChannel(Ball,
		table)

	table <- Ball
	time.Sleep(1 * time.Second)
	KlaudiasGoTrace.ReceiveFromChannel(<-table,
		table)
	KlaudiasGoTrace.EndTrace()

}

//func randStr(n int) string {
//	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}

func playerString(table chan string, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	for {
		ball := <-table
		KlaudiasGoTrace.ReceiveFromChannel(ball,
			table)

		n := (rand.Intn(100-50) + 50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		KlaudiasGoTrace.SendToChannel(ball,
			table)

		//ball += randStr(1)
		table <- ball
	}
	KlaudiasGoTrace.StopGoroutine()

}
