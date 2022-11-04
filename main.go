package main

import (
	"awesomeProject/04-exchange/fanout"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Press 1 to start publisher, 2 to start consumer")
	var startType string
	fmt.Scanf("%s", &startType)
	exchangeType := "direct"
	exchangeName := "router"
	if startType == "1" {
		ch := fanout.InitPublisher(exchangeType, exchangeName)
		if ch == nil {
			return
		}
		fmt.Println("please input a message and a bindingKey in the following format-> key:msg")
		for {
			var msg string
			fmt.Scanf("%s", &msg)
			if msg != "" {
				arr := strings.Split(msg, ":")
				if len(arr) < 2 {
					continue
				}
				fmt.Println("send-->", arr[0], arr[1], len(arr), msg)
				fanout.PublishMessage(arr[1], ch, exchangeName, arr[0])
			}
		}
	} else {
		fmt.Println("please enter a consumer name: ")
		var name string
		for {
			fmt.Scanf("%s", &name)
			if name != "" {
				break
			}
		}
		fmt.Println("please enter a binding key: ")
		var bindingKey string
		for {
			fmt.Scanf("%s", &bindingKey)
			if bindingKey != "" {
				break
			}
		}

		forever := make(chan interface{})
		go fanout.InitConsumer(exchangeType, exchangeName, name, bindingKey)
		<-forever
	}
}
