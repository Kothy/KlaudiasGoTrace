package main

import (
	"math/rand"
	"time"
)

func main() {
	var Ball int
	table := make(chan int)

	go player(table)
	go player(table)

	time.Sleep(20 * time.Millisecond)
	table <- Ball
	time.Sleep(1 * time.Second)
	<-table
}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		n := (rand.Intn(100-50) + 50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		table <- ball
	}
}
