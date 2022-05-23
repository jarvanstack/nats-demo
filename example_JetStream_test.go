package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dengjiawen8955/go_utils/logger"
	nats "github.com/nats-io/nats.go"
)

//注意: 使用 js 之前 nats-server 需要启动 js 模式
var js nats.JetStreamContext

func init() {
	// 初始化 js 上下文
	//256 是一次最大的未完成的异步发送
	js1, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		logger.Fatal(err)
	}
	js = js1
}

func Test_JetStream_DeleteStream(t *testing.T) {
	err := js.DeleteStream("ORDERS")
	if err != nil {
		logger.Fatal(err)
	}
}

//Stream管理,创建流
func Test_JetStream_AddStream(t *testing.T) {
	js.Publish("ORDERS.scratch", []byte("hello"))

	//创建流
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})
	if err != nil {
		logger.Fatal(err)
	}
}
func Test_JetStream_UpdateStream(t *testing.T) {
	//修改流
	_, err := js.UpdateStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"}, //如果这个不是增量更新就需要配置所有,不然会用默认值覆盖
		MaxAge:   time.Second * 10,     //设置10秒过期,然后我们过10秒看看数据还在不在
	})
	if err != nil {
		logger.Fatal(err)
	}
}
func Test_JetStream_AddConsumer(t *testing.T) {
	_, err := js.AddConsumer("ORDERS", &nats.ConsumerConfig{
		Durable:   "MONITOR",
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		logger.Fatal(err)
	}
	//consumer in pull mode requires explicit ack policy
}
func Test_JetStream_DeleteConsumer(t *testing.T) {
	err := js.DeleteConsumer("ORDERS", "MONITOR")
	if err != nil {
		logger.Fatal(err)
	}
}

//-------------
//同步推送
func Test_JetStream_Publish(t *testing.T) {
	js.Publish("ORDERS.scratch", []byte(fmt.Sprintf("hello")))
}

//异步推送
func Test_JetStream_PublishAsync(t *testing.T) {
	for i := 0; i < 5; i++ {
		js.PublishAsync("ORDERS.scratch", []byte(fmt.Sprintf("hello%d", i)))
	}
	select {
	case <-js.PublishAsyncComplete():
		fmt.Println("Publish JetStream Message Success")
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}
}

//消息队列,异步接收
//接收不到消息该死的
func Test_JetStream_Subscribe_Async(t *testing.T) {
	_, err := js.Subscribe("ORDERS.*", func(m *nats.Msg) {
		fmt.Printf("Received a JetStream message: %s\n", string(m.Data))
		m.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
	//防止进程退出收不到消息
	time.Sleep(time.Millisecond * 100)
}

func Test_JetStream_Subscribe_Sync(t *testing.T) {
	sub, _ := js.SubscribeSync("ORDERS.*", nats.Durable("MONITOR"), nats.MaxDeliver(3))
	m, err := sub.NextMsg(time.Second)
	m.Ack()
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("Received a JetStream message: %s\n", string(m.Data))
}
func Test_JetStream_PullSubscribe(t *testing.T) {
	// sub, err := js.PullSubscribe("ORDERS.*", "MONITOR")
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// msgs, err := sub.Fetch(10)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// for _, msg := range msgs {
	// 	fmt.Printf("Received a JetStream message: %s\n", string(msg.Data))
	// }
	//nats: cannot pull subscribe to push based consumer
}
