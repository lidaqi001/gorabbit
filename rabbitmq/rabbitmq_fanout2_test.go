package rabbitmq

import (
	"fmt"
	"github.com/golang-module/carbon"
	"gorabbit/rabbitmq/broker"
	"testing"
)

/**
一对多（广播消息2）
实现：使用topic exchange，保证topic一致，生成随机队列（或自定义）
*/

const TOPIC = "FANOUT"

// 生产者
func TestFanout2Produce(t *testing.T) {
	b := newFanoutBroker2()
	err := b.Publish(
		TOPIC,
		&broker.Message{
			Header: nil,
			Body:   []byte("fanout message"),
		})
	fmt.Println("publish message:", err)
}

// 消费者1
func TestFanout2Consume1(t *testing.T) {
	consumeFanout2("TestConsume1:")
	select {}
}

// 消费者2
func TestFanout2Consume2(t *testing.T) {
	consumeFanout2("TestConsume2:")
	select {}
}

func newFanoutBroker2() broker.Broker {
	DefaultRabbitURL = HOST
	b := NewBroker(
		// 持久化exchange
		//DurableExchange(),
		// 订阅时创建一个持久化队列
		PrefetchGlobal(),
		// 广播时新建一个exchange
		ExchangeName("server_event2"),
	)
	b.Init()
	if err := b.Connect(); err != nil {
		fmt.Printf("cant conect to broker, skip: %v", err)
	}
	return b
}

func consumeFanout2(flag string) {
	b := newFanoutBroker2()
	_, err := b.Subscribe(
		TOPIC,
		func(e broker.Event) error {
			fmt.Println(flag, ", msg:", e.Message(), ",", carbon.Now().ToDateTimeString())
			return nil
		},
		// 持久化队列
		//DurableQueue(),
		// 生成随机队列名
		//broker.Queue(flag),
	)
	fmt.Println(err)
}
