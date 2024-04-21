package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

//
// Common RPC request/reply definitions
//

// 指的是put的args，GetArgs指的是Get的Args
type PutArgs struct {
	Key   string
	Value string
}

// 除了声明时带上了它，其他时候没有用到这个struct
type PutReply struct {
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value string
}

//
// Client
//

func connect() *rpc.Client {
	// Creates a TCP connection to the server.
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

func get(key string) string {
	client := connect()
	args := GetArgs{key} //给stub(KV.Get)传递的参数
	reply := GetReply{}  //用来接收返回值的变量
	// client.Call 方法接受三个参数：RPC方法名 "KV.Get"，参数地址 &args 和回复地址 &reply。如果调用失败，err 被设为对应的错误信息。
	err := client.Call("KV.Get", &args, &reply) //rpc调用
	if err != nil {
		log.Fatal("error:", err)
	}
	client.Close()
	return reply.Value
}

func put(key string, val string) {
	client := connect()
	args := PutArgs{"subject", "6.824"}
	reply := PutReply{}
	// Call() asks the RPC library to perform the call
	// Library marshalls（编码、编组、编集，数据打包） args, sends request, waits, unmarshalls reply
	err := client.Call("KV.Put", &args, &reply)
	if err != nil {
		log.Fatal("error:", err)
	}
	client.Close()
}

//
// Server
//

type KV struct {
	mu   sync.Mutex
	data map[string]string
}

func server() {
	// Go requires server to declare an object with methods as RPC handlers
	kv := &KV{data: map[string]string{}}
	rpcs := rpc.NewServer()

	// Server then registers that object with the RPC library
	rpcs.Register(kv)

	// 监听本地TCP端口1234，reads each request
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	// creates a new goroutine for this request
	go func() {
		for {
			// 循环接受新的连接.
			// 通过监听器 l 等待并接受新的连接。返回的 conn 是指向新连接的对象 
			// Server accepts TCP connections, gives them to RPC library
			conn, err := l.Accept()
			if err == nil {
				// 对于每个新连接，再启动一个协程来处理RPC请求。
				// 对于每个接受的连接，使用 rpcs 服务器来处理该连接的RPC请求。每个连接都在自己的协程中独立处理。
				// 专门用来处理通过 conn 表示的这个新连接的RPC请求。rpcs 是一个RPC服务器实例，ServeConn 方法用于在给定的连接上服务RPC请求。
				go rpcs.ServeConn(conn)
			} else {
				break
			}
		}
		// 当退出循环后，关闭监听器。
		l.Close()
	}()
}

// what will rpc do ?

//unmarshalls request
// looks up the named object (in table create by Register())
// calls the object's named method (dispatch)
// marshalls reply
// writes reply on TCP connection

// The server's Get() and Put() handlers
// Must lock, since RPC library creates a new goroutine for each request
func (kv *KV) Get(args *GetArgs, reply *GetReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	reply.Value = kv.data[args.Key]

	return nil
}

func (kv *KV) Put(args *PutArgs, reply *PutReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[args.Key] = args.Value

	return nil
}

//
// main
//

func main() {
	server()

	put("subject", "6.5840")
	fmt.Printf("Put(subject, 6.5840) done\n")
	fmt.Printf("get(subject) -> %s\n", get("subject"))
}


// Marshalling: format data into packets

// Simplest failure-handling scheme: "best-effort RPC"
//   Call() waits for response for a while
//   If none arrives, re-send the request
//   Do this a few times
//   Then give up and return an error