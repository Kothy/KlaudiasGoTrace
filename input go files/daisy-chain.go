package main

func f(left, right chan int) {
	val := <-right
	val += 1
	left <- val
}

func main() {
	const n = 20
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	go func(c chan int) { c <- 1 }(right)
	<-leftmost
}
