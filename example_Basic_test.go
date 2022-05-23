package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

var nc *nats.Conn

//初始化连接
func init() {
	// nc, _ = nats.Connect(nats.DefaultURL)
	nc, _ = nats.Connect("nats://127.0.0.1:4223")

	//下面只需要调用一个就行了
	// nc.Close()
	// nc.Drain()

	//退订
	// sub.Unsubscribe()
	//所有的 sub 接收者进入 drain 状态, pub 将不能发送任何数据
	// sub.Drain()
	//关闭连接
	// sub.Close()
}

//订阅发布模式,异步接收
func Test_AsyncSubscriber(t *testing.T) {
	//接收数据
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	//模拟发送数据
	go func() {
		nc.Publish("foo", []byte("hello  nats"))
	}()
	//防止进程退出收不到消息
	time.Sleep(time.Millisecond * 100)
}

//订阅发布模式,同步接收
func Test_SyncSubscriber(t *testing.T) {
	// 订阅主题
	sub, err := nc.SubscribeSync("foo")
	if err != nil {
		log.Fatal(err)
	}
	//模拟发送数据
	go func() {
		nc.Publish("foo", []byte("hello  nats 2"))
	}()
	// 阻塞同步获取消息,这个方法会阻塞,直到获取一个消息或者超时
	msg, err := sub.NextMsg(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received a message: %s\n", string(msg.Data))
}

// 队列 模式， 队列为queue, test2 发向所有队列， 只是只是订阅的代码不同而已
//队列模式只有一个接收者能成功接收,用于消息队列, pub,sub 用于广播
func Test_QueueSubscribe(t *testing.T) {
	//接收数据
	//queue是队列组的名称
	_, _ = nc.QueueSubscribe("foo", "queue", func(msg *nats.Msg) {
		fmt.Printf("Queue a message: %s\n", string(msg.Data))
	})
	//模拟发送数据
	go func() {
		nc.Publish("foo", []byte("hello  nats"))
	}()
	//防止进程退出收不到消息
	time.Sleep(time.Millisecond * 100)
}

//订阅发布模式, 使用管道接收
func Test_ChanSubscribe(t *testing.T) {
	//模拟发送数据
	go func() {
		nc.Publish("foo", []byte("hello  nats"))
	}()
	//接收数据
	ch := make(chan *nats.Msg, 64)
	_, _ = nc.ChanSubscribe("foo", ch)
	msg := <-ch
	fmt.Printf("Received a message: %s\n", string(msg.Data))
}

//请求返回模式,
func Test_Request(t *testing.T) {
	//模拟发送数据
	go func() {
		msg, _ := nc.Request("help", []byte("help me"), 10*time.Millisecond)
		fmt.Printf("Received a message2: %s\n", string(msg.Data))
	}()
	//接收数据并返回
	nc.Subscribe("help", func(m *nats.Msg) {
		fmt.Printf("Received a message1: %s\n", string(m.Data))
		nc.Publish(m.Reply, []byte("I can help!"))
	})
	//防止进程退出收不到消息
	time.Sleep(time.Millisecond * 100)
}
