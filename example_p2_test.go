package main

import (
	"fmt"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

func Test_Subscribe2(t *testing.T) {
	// 连接Nats服务器
	nc, _ := nats.Connect("nats://127.0.0.1:4222")

	// 发布-订阅 模式，异步订阅 test1
	_, _ = nc.Subscribe("test1", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	// 队列 模式，订阅 test2， 队列为queue, test2 发向所有队列，同一队列只有一个能收到消息
	_, _ = nc.QueueSubscribe("test2", "queue", func(msg *nats.Msg) {
		fmt.Printf("Queue a message: %s\n", string(msg.Data))
	})

	// 请求-响应， 响应 test3 消息。
	_, _ = nc.Subscribe("test3", func(m *nats.Msg) {
		fmt.Printf("Reply a message: %s\n", string(m.Data))
		_ = nc.Publish(m.Reply, []byte("I can help!!"))
	})

	// 持续发送不需要关闭
	// _ = nc.Drain()

	// 关闭连接
	//nc.Close()

	// 阻止进程结束而收不到消息
	time.Sleep(5 * time.Second)
}

func TestPub2(t *testing.T) {
	// 连接Nats服务器
	nc, _ := nats.Connect("nats://127.0.0.1:4222")

	// 发布-订阅 模式，向 test1 发布一个 `Hello World` 数据
	_ = nc.Publish("test1", []byte("Hello World"))

	// 队列 模式，发布是一样的，只是订阅不同，向 test2 发布一个 `Hello zngw` 数据
	_ = nc.Publish("test2", []byte("Hello zngw"))

	// 请求-响应， 向 test3 发布一个 `help me` 请求数据，设置超时间3秒，如果有多个响应，只接收第一个收到的消息
	msg, err := nc.Request("test3", []byte("help me"), 3*time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("help answer : %s\n", string(msg.Data))
	}

	// 持续发送不需要关闭
	// _ = nc.Drain()

	// 关闭连接
	//nc.Close()
	// 阻止进程结束而收不到消息
	time.Sleep(5 * time.Second)
}
