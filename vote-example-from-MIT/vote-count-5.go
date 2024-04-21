package main

import (
	"math/rand"
	"time"
)

func main() {
	// 注意，main里面不能住能kill一个goroutine，可以send them a var to tell them kill themselves.
	rand.Seed(time.Now().UnixNano())

	// do not need lock here since count and finished are not shared. they are only used once.
	count := 0
	finished := 0
	ch := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			ch <- requestVote()
		}()
	}
	for count < 5 && finished < 10 {
		v := <-ch // 缺点是，first five threads make count reaches to 5， the last five threads will be blocked, hang around. 但是main只要推出，last five threads也会退出。但是让5个线程啥都不做也不大好。
		if v {
			count += 1
		}
		finished += 1
	}
	if count >= 5 {
		println("received 5+ votes!")

	} else {
		println("lost")
	}
}

func requestVote() bool {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return rand.Int()%2 == 0
}
