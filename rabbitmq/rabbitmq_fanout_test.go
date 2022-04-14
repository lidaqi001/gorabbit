package rabbitmq

import (
	"fmt"
	"github.com/golang-module/carbon"
	"gorabbit/rabbitmq/broker"
	"testing"
)

/**
一对多（广播消息1）
实现：使用fanout exchange，生成随机队列（或自定义）
*/

// 生产者
func TestFanoutProduce(t *testing.T) {
	b := newFanoutBroker()
	err := b.PublishFanout(&broker.Message{
		Header: nil,
		Body:   []byte("fanout1 message"),
	})
	fmt.Println("publish message:", err)
}

// 消费者1
func TestFanoutConsume1(t *testing.T) {
	consumeFanout("TestConsume1:")
	select {}
}

// 消费者2
func TestFanoutConsume2(t *testing.T) {
	consumeFanout("TestConsume2:")
	select {}
}

func newFanoutBroker() broker.Broker {
	DefaultRabbitURL = HOST
	b := NewBroker(
		// 持久化exchange
		//DurableExchange(),
		// 订阅时创建一个持久化队列
		PrefetchGlobal(),
		// fanout kind exchange
		FanoutExchange(),
		// 广播时新建一个exchange
		ExchangeName("server_event1"),
	)
	b.Init()
	if err := b.Connect(); err != nil {
		fmt.Printf("cant conect to broker, skip: %v", err)
	}
	return b
}

func consumeFanout(flag string) {
	b := newFanoutBroker()
	_, err := b.SubscribeFanout(
		func(e broker.Event) error {
			fmt.Println(flag, ", msg:", e.Message(), ",", carbon.Now().ToDateTimeString())
			return nil
		},
		// 持久化队列
		//DurableQueue(),
		// 自定义queue标识，为空时随机字符
		//broker.Queue(flag),
	)
	fmt.Println(err)
}
