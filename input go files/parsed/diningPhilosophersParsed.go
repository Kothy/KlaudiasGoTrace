package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Philosopher struct {
	name     string
	fork     chan string
	neighbor *Philosopher
	id       int
	holding  bool
}

func makePhilosopher(name string, neighbor *Philosopher, id int, parentId uint64) *Philosopher {
	phil := &Philosopher{name, make(chan string, 1), neighbor, id, false}
	KlaudiasGoTrace.SendToChannel(
		strconv.Itoa(
			id), phil.fork)

	phil.fork <- strconv.Itoa(id)
	return phil
}

func (phil *Philosopher) think(parentId uint64) {
	fmt.Printf("%v is thinking.\n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (phil *Philosopher) eat(parentId uint64) {
	fmt.Printf("%v is eating.\n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
	fmt.Printf("%v done eating.\n", phil.name)
}

func (phil *Philosopher) getForks(parentId uint64) {
	KlaudiasGoTrace.ReceiveFromChannel(<-phil.fork,
		phil.fork)

	phil.holding = true
	fmt.Printf("%v got his fork.\n", phil.name)
	if !phil.neighbor.holding {
		KlaudiasGoTrace.ReceiveFromChannel(<-phil.neighbor.
			fork,
			phil.
				neighbor.
				fork)

		fmt.Printf("%v got %v's fork.\n", phil.name, phil.neighbor.name)
		fmt.Printf("%v has two forks.\n", phil.name)
		return
	} else {
		KlaudiasGoTrace.SendToChannel(
			strconv.Itoa(
				phil.id), phil.
				fork)

		phil.fork <- strconv.Itoa(phil.id)
		fmt.Printf("%v release his fork.\n", phil.name)
		phil.holding = false
		phil.think(parentId)
		phil.getForks(parentId)
	}
}

func (phil *Philosopher) returnForks(parentId uint64) {
	val1 := strconv.Itoa(phil.id)
	val2 := strconv.Itoa(phil.neighbor.id)
	KlaudiasGoTrace.SendToChannel(
		val1, phil.fork)

	phil.fork <- val1
	fmt.Printf("%v released his fork.\n", phil.name)
	KlaudiasGoTrace.SendToChannel(
		val2, phil.neighbor.
			fork)

	phil.neighbor.fork <- val2
	fmt.Printf("%v released %v's fork.\n", phil.name, phil.neighbor.name)

	phil.holding = false
	phil.neighbor.holding = false
}

func (phil *Philosopher) dine(parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	for {
		phil.think(parentId)
		phil.getForks(parentId)
		phil.eat(parentId)
		phil.returnForks(parentId)
	}
	KlaudiasGoTrace.StopGoroutine()

}

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	names := []string{"Socrates", "Locke", "Descartes", "Newton", "Leibniz"}
	philosophers := make([]*Philosopher, len(names))
	var phil *Philosopher
	for i, name := range names {
		phil = makePhilosopher(name, phil, i, parentId)
		philosophers[i] = phil
	}
	philosophers[0].neighbor = phil

	for _, phil := range philosophers {
		go phil.dine(KlaudiasGoTrace.GetGID())
	}
	time.Sleep(time.Second * 10)
	KlaudiasGoTrace.EndTrace()

}
