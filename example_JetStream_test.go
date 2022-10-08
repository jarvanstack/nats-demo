package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

//注意: 使用 js 之前 nats-server 需要启动 js 模式
var js nats.JetStreamContext

func init() {
	// 初始化 js 上下文
	//256 是一次最大的未完成的异步发送
	js1, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}
	js = js1
}

func TestExample(t *testing.T) {
	// 建立连接
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("连不上 NATS")
	}
	defer nc.Close()

	// 转换为 JetStream 的连接
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("转换为 JetStream 的连接: %v", err)
	}

	// 新增 Stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name: "STREAM_NAME",
		Subjects: []string{
			"foo",
		},
		Storage:   nats.FileStorage,     // 消息储存方式
		Retention: nats.WorkQueuePolicy, // 消息保留策略
		Discard:   nats.DiscardOld,      // 消息丢弃的策略
		// ...
	})
	if err != nil {
		log.Fatalf("建立 Stream 失敗: %v", err)
	}

	// 新增 Consumer
	js.AddConsumer("STREAM_NAME", &nats.ConsumerConfig{
		Durable:   "CONSUMER_NAME",
		AckPolicy: nats.AckExplicitPolicy, // 确认的策略
		// ...
	})

	// 推送消息
	_, _ = js.Publish("foo", []byte("Hello World!"))

	// 模拟 push 模式
	go func() {
		_, err = js.Subscribe("subject", func(msg *nats.Msg) {
			fmt.Println("收到了", string(msg.Data))
		})
		if err != nil {
			log.Fatal("订阅失败", err)
		}
	}()

	// 模拟 pull 模式
	sub, err := js.PullSubscribe("subject", "durable") // pull 模式需要用 Durable
	if err != nil {
		log.Fatalf("订阅失败: %v", err)
	}
	for {
		msgs, err := sub.Fetch(10) // 设置 1 次接受 10 条
		if err != nil {
			log.Fatalf("接收失败: %v", err)
		}
		for _, msg := range msgs {
			fmt.Println("收到了", string(msg.Data))
			msg.Ack()
		}
	}
}

func Test_JetStream_DeleteStream(t *testing.T) {
	err := js.DeleteStream("ORDERS")
	if err != nil {
		log.Fatal(err)
	}
}

//Stream管理,创建流
func Test_JetStream_AddStream(t *testing.T) {
	js.Publish("ORDERS.scratch", []byte("hello"))

	//创建流
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"foo"},
	})
	if err != nil {
		log.Fatal(err)
	}
}
func Test_JetStream_UpdateStream(t *testing.T) {
	//修改流
	_, err := js.UpdateStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"foo"},  //如果这个不是增量更新就需要配置所有,不然会用默认值覆盖
		MaxAge:   time.Second * 10, //设置10秒过期,然后我们过10秒看看数据还在不在
	})
	if err != nil {
		log.Fatal(err)
	}
}
func Test_JetStream_AddConsumer(t *testing.T) {
	_, err := js.AddConsumer("ORDERS", &nats.ConsumerConfig{
		Durable:   "MONITOR",
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}
	//consumer in pull mode requires explicit ack policy
}
func Test_JetStream_DeleteConsumer(t *testing.T) {
	err := js.DeleteConsumer("ORDERS", "MONITOR")
	if err != nil {
		log.Fatal(err)
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
	_, err := js.Subscribe("foo", func(m *nats.Msg) {
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
	sub, _ := js.SubscribeSync("foo", nats.Durable("MONITOR"), nats.MaxDeliver(3))
	m, err := sub.NextMsg(time.Second)
	m.Ack()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received a JetStream message: %s\n", string(m.Data))
}
func Test_JetStream_PullSubscribe(t *testing.T) {
	sub, err := js.PullSubscribe("foo", "durable")
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := sub.Fetch(10)
	if err != nil {
		log.Fatal(err)
	}
	for _, msg := range msgs {
		fmt.Printf("Received a JetStream message: %s\n", string(msg.Data))
	}
	// nats: cannot pull subscribe to push based consumer
}
