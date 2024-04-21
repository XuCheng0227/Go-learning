package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d started\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d completed\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers have completed.")
}

// sync.WaitGroup 是 Go 语言标准库中 sync 包提供的一种同步原语，用于等待一组 goroutine 完成任务，是一个相当常用的数据结构。

// 它常用于等待一批并发操作完成后再继续执行主程序，以确保并发操作完成后再进行后续的操作。

// 在 sync.WaitGroup（等待组）类型中，每个 sync.WaitGroup值在内部维护着一个计数器(此计数的初始默认值为0)，sync.WaitGroup 通过计数器来追踪等待的 goroutine 数量，每个等待的 goroutine 在开始执行之前都会调用 Add 方法来增加计数器。当 goroutine 执行完任务后，应调用 Done 方法来减少计数器。主程序可以调用 Wait 方法来阻塞，直到计数器归零，即所有的 goroutine 都执行完毕。

// sync.WaitGroup 提供了如下几个函数可用：

// func (wg *WaitGroup) Add(delta int) {...}
// func (wg *WaitGroup) Done() {...}
// func (wg *WaitGroup) Wait() {...}
// Add(delta int)：增加计数器的值。 delta 表示要增加的计数器值，可以为正数或负数。
// Done()：减少计数器的值，相当于执行完一个任务。
// Wait()：阻塞主程序，直到计数器归零。一般在主程序中调用，用于等待所有任务完成。
