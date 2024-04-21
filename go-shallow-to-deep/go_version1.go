package main

import (
	"fmt"
	"sync"
	"time"
)

var cookies int
var m sync.Mutex
var wg sync.WaitGroup
var line_ch, final_ch chan string

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
		final_ch <- name
	} else {
		fmt.Printf("%s found no cookies...Leave and wait in line again.\n", name)
		line_ch <- name // 没吃到就重新排队
	}
	return
}

func makeCookies() {
	m.Lock()
	defer m.Unlock()
	// 这一版代码中，makeCookies不需要wg.Done()。
	//1. 如果defer一个go func()，chef会在main里不断makeCookies，
	// 但是wg.Add(1)只加了一次，makeCookies里面却wg.Done了好几次。程序会panic
	//2. 因为main里面chef有一个for循环去按照一定的时间去检查有没有cookies。
	// 如果还不断wg.Add(1)，那么main结束不了，for循环也不会退出。

	if cookies <= 0 {
		fmt.Printf("The chef is seated.\n")
		fmt.Printf("Making cookies at this particular table.\n")
		time.Sleep(time.Second)
		cookies += 5
		fmt.Printf("%d cookies are available for new customers.\n", cookies)
	} else {
		fmt.Printf("Cookies are still enough, do nothing.\n")
	}

}

func main() {
	cookies = 5
	customers := [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
	line_ch = make(chan string, len(customers))
	final_ch = make(chan string, len(customers))

	for _, individual := range customers {
		line_ch <- individual
	}

	// customers eat cookies
	go func() {
		for person := range line_ch {
			wg.Add(1)
			go eatCookies(person)
		}
	}()

	// chef make cookies
	wg.Add(1)
	go func() {
		for {
			time.Sleep(2 * time.Second) // 每隔这么一段时间就进去看看cookies还够不够，不够就继续做cookies。
			makeCookies()
			fmt.Printf("final_ch: %d\n", len(final_ch))
			if len(customers) == len(final_ch) {
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

//缺点：有人发现cookie没了，应该是叫厨师先做，而不应该重新去排队。厨师也不是隔一段时间过来看下cookie有没有剩余，是否要make cookies。
//问题改进：如果有人发现没有cookies，应该在旁边等着。等厨师做好cookies，然后通知大家cookie出锅，可以继续排队了。
