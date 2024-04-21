package main

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

var cookies int
var m sync.Mutex
var cond = sync.NewCond(&m)

var line_ch, final_ch chan string

func eatCookies(name string) {
	m.Lock()
	defer m.Unlock()
	fmt.Printf("%s enter into the restaurant\n", name)

	for cookies < 1 {
		fmt.Printf("%s leaves the restaurant and wait for the complete signal from the chef\n", name)
		cond.Wait() // 此处相当于调用了m.Unlock()。待后续收到信号后再重新m.Lock()。
		fmt.Printf("New cookies are made and %s returns.\n", name)
	}

	fmt.Printf("Eating cookies at this particular table...\n")
	time.Sleep(time.Second)
	cookies -= 1
	fmt.Printf("%s has finished. Prepare to leave\n", name)

	if cookies < 1 {
		cond.Broadcast() // 不能用cond.signal()， signal不是通知特定的goroutine，而是只通知一个goroutine。
		fmt.Printf("Notify all so that chef can receice the sign to make cookies")
	}
	fmt.Printf("%s is leaving\n", name)
	return
}

func makeCookies() {
	m.Lock()
	defer m.Unlock()
	fmt.Printf("The chef enter into the restaurant\n")

	for cookies > 0 {
		fmt.Printf("Cookies are still enough, do nothing and wait for the sign.\n")
		cond.Wait()
		fmt.Printf("The chef has received the sign.\n")
	}

	fmt.Printf("The chef is making cookies at this particular table.\n")
	time.Sleep(time.Second)
	cookies += 5
	cond.Broadcast()
	fmt.Printf("%d cookies are available for new customers, the chef is leaving.\n", cookies)
}

func main() {
	cookies = 5
	customers := [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	go func() { // The table is available forever.
		for i := 0; i < math.MaxInt; i++ {
			for _, person := range customers {
				real_person := person + strconv.Itoa(i)
				go eatCookies(real_person)
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	// chef make cookies
	go func() {
		for {
			makeCookies()
		}
	}()

	// 缓冲的通道来阻塞主函数的执行，直到所有的goroutine执行完毕。主函数会一直阻塞在<-end这一行，直到有goroutine向end通道发送数据，才会继续执行并结束程序。这样可以确保所有的goroutine都有足够的时间运行，并输出结果。
	end := make(chan bool)
	<-end

	// 如果不加这一行，或者不加缓冲的通道阻塞主函数的执行，主函数在创建goroutine后立即结束，没有等待goroutine的执行。当主函数结束时，程序也会随之结束，所有的goroutine便没有足够的时间运行。
	//time.Sleep(time.Second * 10)

}
