package fanout

import (
	"awesomeProject/common/mqUtils"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitConsumer(exchangeType, exchangeName, consumerName, bindingKey string) {
	//1. 建立mq链接
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", mqUtils.MQ_USER, mqUtils.MQ_PWD, mqUtils.MQ_ADDR)
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("mq connection error: ", err)
		return
	}

	//2. 创建信道
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("channel creation error: ", err)
		return
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
		return
	}

	q := CreateTemporaryQueue(ch, exchangeName, bindingKey)
	if q == nil {
		fmt.Println("create queue failure")
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("consume error: ", err)
		return
	}
	forever := make(chan interface{})
	go func() {
		for msg := range msgs {
			fmt.Println("receive ---> ", string(msg.Body))
		}
	}()
	<-forever
}

func CreateTemporaryQueue(ch *amqp.Channel, exchangeName string, bindingKey string) *amqp.Queue {
	if ch == nil {
		fmt.Println("channel is nil")
		return nil
	}

	//1. 声明临时队列
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		fmt.Println("temporary queue declaration error: ", err)
		return nil
	}

	//2. 队列和交换机绑定
	err = ch.QueueBind(
		q.Name,
		bindingKey,
		exchangeName,
		false,
		nil)
	if err != nil {
		fmt.Println("queue binding error: ", err)
		return nil
	}
	return &q
}
