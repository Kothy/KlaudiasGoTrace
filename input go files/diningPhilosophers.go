package main

import (
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

func makePhilosopher(name string, neighbor *Philosopher, id int) *Philosopher {
	phil := &Philosopher{name, make(chan string, 1), neighbor, id, false}
	phil.fork <- strconv.Itoa(id)
	return phil
}

func (phil *Philosopher) think() {
	fmt.Printf("%v is thinking.\n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (phil *Philosopher) eat() {
	fmt.Printf("%v is eating.\n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
	fmt.Printf("%v done eating.\n", phil.name)
}

func (phil *Philosopher) getForks() {
	<-phil.fork
	phil.holding = true
	fmt.Printf("%v got his fork.\n", phil.name)
	if !phil.neighbor.holding {
		<-phil.neighbor.fork
		fmt.Printf("%v got %v's fork.\n", phil.name, phil.neighbor.name)
		fmt.Printf("%v has two forks.\n", phil.name)
		return
	} else {
		phil.fork <- strconv.Itoa(phil.id)
		fmt.Printf("%v release his fork.\n", phil.name)
		phil.holding = false
		phil.think()
		phil.getForks()
	}
}

func (phil *Philosopher) returnForks() {
	val1 := strconv.Itoa(phil.id)
	val2 := strconv.Itoa(phil.neighbor.id)
	phil.fork <- val1
	fmt.Printf("%v released his fork.\n", phil.name)

	phil.neighbor.fork <- val2
	fmt.Printf("%v released %v's fork.\n", phil.name, phil.neighbor.name)

	phil.holding = false
	phil.neighbor.holding = false
}

func (phil *Philosopher) dine() {
	for {
		phil.think()
		phil.getForks()
		phil.eat()
		phil.returnForks()
	}
}

func main() {
	names := []string{"Socrates", "Locke", "Descartes", "Newton", "Leibniz"}
	philosophers := make([]*Philosopher, len(names))
	var phil *Philosopher
	for i, name := range names {
		phil = makePhilosopher(name, phil, i)
		philosophers[i] = phil
	}
	philosophers[0].neighbor = phil

	for _, phil := range philosophers {
		go phil.dine()
	}
	time.Sleep(time.Second * 10)
}
