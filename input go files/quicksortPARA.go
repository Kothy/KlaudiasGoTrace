package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var (
	size        = 750
	granularity = 2
)

func pivot(pole []int) int {
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

func quickSort(pole []int) {
	if len(pole) > granularity {
		done := make(chan bool)
		go cquickSort(pole, done)
		<-done
	} else {
		squickSort(pole)
	}
}

func cquickSort(pole []int, done chan bool) {
	if len(pole) <= 1 {
		done <- true
	} else if len(pole) < granularity {
		quickSort(pole)
		done <- true
	} else {
		index := pivot(pole)
		left, right := make(chan bool), make(chan bool)
		go cquickSort(pole[:(index-1)], left)
		go cquickSort(pole[index:], right)
		l := <-left
		r := <-right
		done <- (l && r)
	}
}

func squickSort(pole []int) {
	if len(pole) > 1 {
		index := pivot(pole)
		squickSort(pole[:(index - 1)])
		squickSort(pole[index:])
	}
}

func main() {

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
	quickSort(array)
	t1 := time.Now()

	fmt.Println(array)
	fmt.Println("Sorted in", t1.Sub(t0).Nanoseconds()/1000000, "milliseconds.")

}
