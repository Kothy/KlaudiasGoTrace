package main

import (
	"math/rand"
	"time"
)

func main() {
	var Ball string
	Ball = "Ball"
	table := make(chan string)

	go playerString(table)
	go playerString(table)

	time.Sleep(20 * time.Millisecond)
	table <- Ball
	time.Sleep(1 * time.Second)
	<-table
}

//func randStr(n int) string {
//	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}

func playerString(table chan string) {
	for {
		ball := <-table
		n := (rand.Intn(100-50) + 50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		//ball += randStr(1)
		table <- ball
	}
}
