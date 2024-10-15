package mq

import (
	"IM/config"
	"IM/model"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

const (
	MessageQueue        = "message.queue"
	MessageExchangeType = "direct"
	MessageRouteKey     = "message.route.key"
	MessageExchangeName = "message.exchange.name"
)

var RabbitMQ *RabbitMQConn

type RabbitMQConn struct {
	Conn         *amqp091.Connection
	Producer     *amqp091.Channel
	Consumer     <-chan amqp091.Delivery
	ExchangeName string
	RouteKey     string
}

// 初始化RabbitMQ
func InitRabbitMQ(config *config.RabbitMQConfig) (err error) {
	conn, err := amqp091.Dial(config.Url)
	if err != nil {
		return err
	}
	producer, err := NewProducer(conn)
	if err != nil {
		return err
	}

	consumer, err := NewConsumer(conn, handleDelivery)
	RabbitMQ = &RabbitMQConn{
		Conn:         conn,
		Producer:     producer,
		Consumer:     consumer,
		ExchangeName: MessageExchangeName,
		RouteKey:     MessageRouteKey,
	}

	return nil
}

// 注册生产者
func NewProducer(conn *amqp091.Connection) (*amqp091.Channel, error) {
	// 1.创建通道
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 2.声明交换机
	err = ch.ExchangeDeclare(
		MessageExchangeName, // 交换机名称
		MessageExchangeType, // 交换机类型
		true,                // 是否持久化
		false,               // 是否自动删除
		false,               // 是否内部使用
		false,               // 是否等待服务器响应
		nil,                 // 其他属性
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

// 注册消费者
func NewConsumer(conn *amqp091.Connection, handleDelivery func(amqp091.Delivery)) (<-chan amqp091.Delivery, error) {

	// 1.创建通道
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 2.声明交换机
	err = ch.ExchangeDeclare(
		MessageExchangeName, // 交换机名称
		MessageExchangeType, // 交换机类型
		true,                // 是否持久化
		false,               // 是否自动删除
		false,               // 是否内部使用
		false,               // 是否等待服务器响应
		nil,                 // 其他属性
	)
	if err != nil {
		return nil, err
	}

	// 3.声明队列
	q, err := ch.QueueDeclare(MessageQueue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// 4.绑定队列到交换机
	err = ch.QueueBind(q.Name, MessageRouteKey, MessageExchangeName, false, nil)
	if err != nil {
		return nil, err
	}

	// 5.订阅消息
	consumer, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// 启动一个goroutine来处理消息
	go func() {
		for d := range consumer {
			handleDelivery(d)
		}
	}()
	return consumer, nil
}

// 消息处理函数
func handleDelivery(d amqp091.Delivery) {
	//fmt.Printf("Received message: %s\n", d.Body)
	messages := make([]model.Message, 0)

	err := json.Unmarshal(d.Body, &messages)
	if err != nil {
		fmt.Println("handleDelivery.json.Unmarshal Error:", err)
		err = d.Nack(false, false)
		if err != nil {
			fmt.Println("handleDelivery.Nack Error:", err)
			return
		}
		return
	}
	//fmt.Println("handleDelivery.message:", messages)
	// 保存消息到数据库
	//err = model.CreateMessages(messages)
	//if err != nil {
	//	fmt.Println("handleDelivery.model.CreateMessages Error:", err)
	//	return
	//}
	// 确认消息
	if err = d.Ack(false); err != nil {
		log.Printf("Failed to acknowledge delivery: %v", err)
	}
	fmt.Println("消费者处理完成")
}
