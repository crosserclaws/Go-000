package main

import (
	"log"
	"time"
	"week06/pkg/window"
)

func main() {
	log.Println("Start.")

	opt := &window.SlidingWindowOption{Capacity: 10, Unit: time.Millisecond * time.Duration(500)}
	success := window.NewSlidingWindow(opt)

	for i := 1; i <= 15; i++ {
		if i <= 10 {
			add(success, 2)
		}
		log.Printf("Window: t=%02v, start=%v, success sum of window=%v", i, success.GetStartTime(), success.Sum())
		sleep(time.Second)
	}
	log.Println("End")
}

func add(w *window.SlidingWindow, n int64) {
	log.Println("+", n)
	w.Add(n)
}

func sleep(d time.Duration) {
	log.Println("Sleep:", d)
	time.Sleep((d))
}
