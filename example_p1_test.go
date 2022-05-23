package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

//先开 sub 再开 pub 消息是无法堆积的
//异步阻塞
func Test_NextMsg(t *testing.T) {
	// 订阅主题
	sub, err := nc.SubscribeSync("foo")
	if err != nil {
		log.Fatal(err)
	}
	// 等待下一个消息消息，并且设置超时时间
	msg, err := sub.NextMsg(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Reply: %s", msg.Data)

}

//异步订阅
//这个总是测试失败???
func Test_Subscribe(t *testing.T) {
	// 异步订阅
	if _, err := nc.Subscribe("foo", func(msg *nats.Msg) {
		fmt.Printf("msg.Data: %s\n", msg.Data)
	}); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
}

func Test_Publish(t *testing.T) {
	// Simple Publisher
	for i := 0; i < 10; i++ {
		err := nc.Publish("foo", []byte(fmt.Sprintf("Hello World%d", i)))
		if err != nil {
			log.Fatal(err)
		}
	}
	//这个是异步发送的,可能还没有完全发出去,太快了,暂停下
	time.Sleep(100 * time.Microsecond)
}
