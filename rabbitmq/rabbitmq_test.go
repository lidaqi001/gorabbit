package rabbitmq

import (
	"fmt"
	"github.com/golang-module/carbon"
	"gorabbit/rabbitmq/broker"
	"testing"
)

/**
一对一（消息队列）
实现：使用topic类型
*/

const HOST = "amqp://admin:123456@127.0.0.1:5672"
const QUEUE = "test_queue"

// 生产者
func TestProduce(t *testing.T) {
	b := newBroker()
	err := b.Publish(QUEUE, &broker.Message{
		Header: nil,
		Body:   []byte("test message"),
	})
	fmt.Println("publish message:", err)
}

// 消费者1
func TestConsume1(t *testing.T) {
	consume(QUEUE, "TestConsume1:")
	select {}
}

// 消费者2
func TestConsume2(t *testing.T) {
	consume(QUEUE, "TestConsume2:")
	select {}
}

func newBroker() broker.Broker {
	DefaultRabbitURL = HOST
	b := NewBroker(
		// 持久化exchange
		DurableExchange(),
		// 订阅时创建一个持久化队列
		PrefetchGlobal(),
	)
	b.Init()
	if err := b.Connect(); err != nil {
		fmt.Printf("cant conect to broker, skip: %v", err)
	}
	return b
}

func consume(queue, flag string) {
	b := newBroker()
	_, err := b.Subscribe(
		queue,
		func(e broker.Event) error {
			fmt.Println(flag, ", msg:", e.Message(), ",", carbon.Now().ToDateTimeString())
			return nil
		},
		// 持久化队列
		DurableQueue(),
		broker.Queue(queue),
	)
	fmt.Println(err)
}
