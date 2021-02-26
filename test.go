package KlaudiasGoTrace

import (
	"time"
)

func producer(ch chan int, d time.Duration, parentId uint64) {
	StartGoroutine(parentId)
	for i := 0; i < 10; i++ {
		go myGo(GetGID())
	}

	for i := 0; i < 10; i++ {
		SendToChannel(i, ch)
		ch <- i
		time.Sleep(d)
	}
	StopGoroutine()
}

func myGo(parentId uint64) {
	StartGoroutine(parentId)
	time.Sleep(10 * time.Millisecond)
	StopGoroutine()
}

func reader(out chan int, parentId uint64) {
	StartGoroutine(parentId)
	for i := 0; i < 20; i++ {
		ReceiveFromChannel(<-out, out)
	}
	StopGoroutine()
}

//func main() {
//	KlaudiasGoTrace.StartTrace()
//
//	ch := make(chan int)
//	out := make(chan int)
//
//	go producer(ch, 10 * time.Millisecond, KlaudiasGoTrace.GetGID())
//	go producer(ch, 25 * time.Millisecond, KlaudiasGoTrace.GetGID())
//	go reader(out, KlaudiasGoTrace.GetGID())
//
//	for i := 0; i < 20; i++ {
//		i := <-ch
//		KlaudiasGoTrace.ReceiveFromChannel(i, ch)
//		KlaudiasGoTrace.SendToChannel(i, out)
//		out <- i
//	}
//
//	KlaudiasGoTrace.EndTrace()
//}
