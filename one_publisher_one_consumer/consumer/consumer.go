package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	username = "username"
	password = "secret"
	queue    = "myqueue"
)

type Data struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
}

func main() {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, "localhost", 5672))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	delivery, err := ch.Consume(
		q.Name,
		"my consumer",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	for msg := range delivery {
		fmt.Println("Message is recieved", string(msg.Body))
	}
}
