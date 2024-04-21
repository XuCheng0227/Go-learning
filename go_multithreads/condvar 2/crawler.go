package main

import (
	"fmt"
	"sync"
)

//
// Several solutions to the crawler exercise from the Go tutorial
// https://tour.golang.org/concurrency/10
//

//
// Serial crawler
//

func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return
	}
	fetched[url] = true
	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	for _, u := range urls {
		Serial(u, fetcher, fetched)
	}
	return
}

//
// Concurrent crawler with shared state and Mutex
//

// the map needs to be protected by mutex. Map itself is not thread safe. 
type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func (fs *fetchState) testAndSet(url string) bool {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	r := fs.fetched[url]
	fs.fetched[url] = true
	return r
}

func ConcurrentMutex(url string, fetcher Fetcher, fs *fetchState) {
	if fs.testAndSet(url) { // 某个URL已经被fetched了，就直接跳过。
		return
	}
	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	var done sync.WaitGroup // a primitive to keep track of how many active threads you still have and when to terminate.
	for _, u := range urls {
		done.Add(1)
		go func(u string) {
			defer done.Done()
			ConcurrentMutex(u, fetcher, fs)
		}(u)
	}
	done.Wait()
	return
}

func makeState() *fetchState {
	return &fetchState{fetched: make(map[string]bool)}
}

//
// Concurrent crawler with channels
//

func worker(url string, ch chan []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url)
	if err != nil {
		ch <- []string{}
	} else {
		ch <- urls
	}
}

func coordinator(ch chan []string, fetcher Fetcher) {
	n := 1 // keep track of how many outstanding workers and coordinator we have.
			//n的初始值赋值为1，这个1指的是coordinator一开始就接受了一个chan中的url。
			// Keeps count of workers in n. Each worker sends exactly one item on channel.
	fetched := make(map[string]bool)
	for urls := range ch {
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n += 1
				go worker(u, ch, fetcher) // worker的作用是把一个URL抓到的URLs放到channel里面。
			}
		}
		n -= 1 //指所有URL都被遍历到了
		if n == 0 {
			break
		}
	}
}

func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() { // 此处必须要用goroutine，否则就deadlock了。
		// sender blocks until the receiver receives!
		ch <- []string{url}
	}()
	coordinator(ch, fetcher)
}

// 上述函数如果去掉go这个关键字，发送数据到channel是一个阻塞性操作，如果没有其他goroutine等待接收这些数据，那么发送操作会一直等待，从而阻塞当前的goroutine。由于在上面的代码中，发送数据(ch <- []string{url})和接收数据(coordinator(ch, fetcher))是在同一个goroutine中顺序执行的，因此发送操作将永远等待，因为它阻塞了goroutine，导致coordinator函数永远不会被调用，从而导致死锁。
// 原始代码使用go func()来启动一个新的goroutine仅用于发送数据，这允许主goroutine继续执行并调用coordinator函数，该函数可以从channel接收数据，从而避免了死锁。因为发送者和接收者分别在不同的goroutines中，这样就可以同时进行发送和接收操作。
//总结：如果发送和接收操作在同一个goroutine中顺序执行，那么发送操作将会阻塞该goroutine，从而阻止接收操作执行，导致死锁。这就是为什么去掉go关键字会导致死锁的原因。


//
// main
//

func main() {
	fmt.Printf("=== Serial===\n")
	Serial("http://golang.org/", fetcher, make(map[string]bool))

	fmt.Printf("=== ConcurrentMutex ===\n")
	ConcurrentMutex("http://golang.org/", fetcher, makeState())

	fmt.Printf("=== ConcurrentChannel ===\n")
	ConcurrentChannel("http://golang.org/", fetcher)
}

//
// Fetcher
//

type Fetcher interface {
	// Fetch returns a slice of URLs found on the page.
	Fetch(url string) (urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

// fakeFetcher实现了Fetch方法，因此实现了Fetcher这个接口
func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("found:   %s\n", url)
		return res.urls, nil
	}
	fmt.Printf("missing: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{  // 注意，这个fakeFetcher类型的fetcher，最终要传到main函数里，和各种方法里面的Fetcher类型的fetcher配对。这样就可以直接调用Fetcher类型的fetcher.fetch()了。
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}