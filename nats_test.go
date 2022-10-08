package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

func TestSubscribe(t *testing.T) {
	// 生产消息
	go func() {
		nc.Publish("foo", []byte("hello nats"))
	}()
	// 消费消息
	_, err := nc.Subscribe("foo", func(msg *nats.Msg) {
		fmt.Printf("收到消息: %s\n", msg.Data)
	})
	if err != nil {
		log.Fatal("NATS 订阅失败")
	}
	select {}

}

func TestRequestReply(t *testing.T) {
	// 消费消息
	nc.Subscribe("help", func(m *nats.Msg) {
		fmt.Printf("收到消息: %s\n", string(m.Data))
		nc.Publish(m.Reply, []byte("I can help!"))
	})

	// 生产消息
	go func() {
		msg, _ := nc.Request("help", []byte("help me"), 100*time.Millisecond)
		fmt.Printf("收到回复: %s\n", string(msg.Data))
	}()

	select {}
}
