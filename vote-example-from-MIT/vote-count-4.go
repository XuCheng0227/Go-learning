package main

import (
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	count := 0
	finished := 0
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	for i := 0; i < 10; i++ {
		go func() {
			vote := requestVote()
			mu.Lock()
			defer mu.Unlock()
			if vote {
				count++
			}
			finished++
			// unlike channels, blocked when send values to channels. cond won't block the thread.
			cond.Broadcast()
		}()
	}

	mu.Lock()
	for count < 5 && finished != 10 {
		cond.Wait() // wait要先做一个unlock的操作。直到收到broadcast的信号才结束wait，继续完成for循环的判断条件。
		// cond.Wait() 会把调用者Caller放入Cond的等待队列中并阻塞，直到被Signal或者Broadcast的方法从等待队列中移除并唤醒。
		// 这种使用 sync.Cond 的模式使得我们可以在条件尚未满足时有效地等待，而不是进行无意义的轮询（busy waiting），一旦条件有可能被满足时，就会通知等待的 goroutines 进行检查，这种方式更加高效。
		// 如果循环条件为真，则调用 cond.Wait()。cond.Wait() 会做两件事：1. 它会自动释放关联的 mu 互斥锁，允许其他 goroutines 获取这个锁来修改 count 和 finished 变量。2. 它会阻塞当前的主 goroutine，直到 cond.Broadcast() 被调用。
		// 一旦 cond.Broadcast() 被任意一个投票 goroutine 调用，主 goroutine 会从阻塞状态被唤醒，自动重新获取 mu 互斥锁，并回到循环的条件检查部分。
	}
	if count >= 5 {
		println("received 5+ votes!")
	} else {
		println("lost")
	}
	mu.Unlock()
}

func requestVote() bool {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return rand.Int()%2 == 0
}
