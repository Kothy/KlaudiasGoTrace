package main

import (
	"KlaudiasGoTrace/KlaudiasGoTrace"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var (
	size        = 750
	granularity = 2
)

func pivot(pole []int, parentId uint64) int {
	i, j, x := 1, len(pole)-1, pole[0]
	for i <= j {
		for i <= j && pole[i] <= x {
			i++
		}
		for j >= i && pole[j] >= x {
			j--
		}
		if i < j {
			pole[i], pole[j] = pole[j], pole[i]
		}
	}
	pole[0], pole[j] = pole[j], pole[0]
	return i
}

func quickSort(pole []int, parentId uint64) {
	if len(pole) > granularity {
		done := make(chan bool)
		go cquickSort(pole, done, KlaudiasGoTrace.GetGID())
		KlaudiasGoTrace.ReceiveFromChannel(<-done,
			done)

	} else {
		squickSort(pole, parentId)
	}
}

func cquickSort(pole []int, done chan bool, parentId uint64) {
	KlaudiasGoTrace.StartGoroutine(parentId)

	if len(pole) <= 1 {
		KlaudiasGoTrace.SendToChannel(
			true, done)

		done <- true
	} else if len(pole) < granularity {
		quickSort(pole, parentId)
		KlaudiasGoTrace.SendToChannel(
			true, done)

		done <- true
	} else {
		index := pivot(pole, parentId)
		left, right := make(chan bool), make(chan bool)
		go cquickSort(pole[:(index-1)], left, KlaudiasGoTrace.GetGID())
		go cquickSort(pole[index:], right, KlaudiasGoTrace.GetGID())
		l := <-left
		KlaudiasGoTrace.ReceiveFromChannel(l, left)

		r := <-right
		KlaudiasGoTrace.ReceiveFromChannel(r, right)
		KlaudiasGoTrace.SendToChannel(
			(l && r), done)

		done <- (l && r)
	}
	KlaudiasGoTrace.StopGoroutine()

}

func squickSort(pole []int, parentId uint64) {
	if len(pole) > 1 {
		index := pivot(pole, parentId)
		squickSort(pole[:(index-1)], parentId)
		squickSort(pole[index:], parentId)
	}
}

func main() {
	parentId := uint64(0)
	KlaudiasGoTrace.Use(parentId)
	KlaudiasGoTrace.StartTrace()

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Number of CPU cores:", runtime.NumCPU())
	array := make([]int, size)
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10000
	for i, _ := range array {
		array[i] = rand.Intn(max-min) + min
	}
	fmt.Println(array)
	t0 := time.Now()
	quickSort(array, parentId)
	t1 := time.Now()

	fmt.Println(array)
	fmt.Println("Sorted in", t1.Sub(t0).Nanoseconds()/1000000, "milliseconds.")
	KlaudiasGoTrace.EndTrace()

}
