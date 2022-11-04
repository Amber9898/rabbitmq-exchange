package fanout

import (
	"awesomeProject/common/mqUtils"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func InitPublisher(exchangeType string, exchangeName string) *amqp.Channel {
	//QUEUE_NAME := fmt.Sprintf("Consumer-%d", time.Now().Unix())
	//durable := false //临时队列

	//1. 建立mq链接
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", mqUtils.MQ_USER, mqUtils.MQ_PWD, mqUtils.MQ_ADDR)
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("mq connection error: ", err)
		return nil
	}

	//2. 创建信道
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("channel creation error: ", err)
		return nil
	}

	//3. 声明交换机
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println("exchange declaration error: ", err)
		return nil
	}
	return ch
}

func PublishMessage(msg string, ch *amqp.Channel, exchangeName string, bindingKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ch.PublishWithContext(
		ctx,
		exchangeName,
		bindingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
}
