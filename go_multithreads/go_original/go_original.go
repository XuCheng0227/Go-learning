package main

import (
	"fmt"
	"sync"
	"time"
)

var cookies int
var m sync.Mutex
var wg sync.WaitGroup

func eatCookies(name string) {
	m.Lock()
	defer func() {
		wg.Done()
		m.Unlock()
	}()

	if cookies >= 1 {
		fmt.Printf("%s is seated.\t", name)
		fmt.Printf("Eating cookies at this particular table...\n")
		time.Sleep(time.Second)
		cookies -= 1
		fmt.Printf("%s has finished. Prepare to leave\n", name)
		fmt.Printf("cookies: %d\n", cookies)
	} else {
		fmt.Printf("%s found no cookies...Leaving.\n", name)
	}
	return
}

func makeCookies() {
	m.Lock()
	defer func() {
		wg.Done()
		m.Unlock()
	}()

	if cookies <= 0 {
		fmt.Printf("The chef is seated.\n")
		fmt.Printf("Making cookies at this particular table.\n")
		fmt.Printf("cookies: %d\n", cookies)
		time.Sleep(time.Second)
		cookies += 5
		fmt.Printf("Cookies are made for new customers.\n")
		fmt.Printf("cookies: %d\n", cookies)
	} else {
		fmt.Printf("cookies: %d\n", cookies)
		fmt.Printf("Cookies are still enough, do nothing.\n")
	}

}

func main() {
	cookies = 1
	customers := [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	for _, individual := range customers {
		wg.Add(1)
		go eatCookies(individual)
	}

	wg.Add(1) // 如果不使用wg，main执行完以后可能直接结束了。因此需要使用wg让main等待所有携程执行完才可以继续往下执行。
	go makeCookies()
	wg.Wait()
}

// 这段代码的问题是每个人只排队在这个桌子上吃或做一次cookie。
// customres发现没有cookie应该让厨师继续做，而不应该直接走。厨师也不应该只检查一次
