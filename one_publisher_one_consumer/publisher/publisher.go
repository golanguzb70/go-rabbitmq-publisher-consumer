package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	// publishing 10 messages. 1 in every second.
	for i := 0; i < 10; i++ {
		message := Data{
			Message: "Data message",
			Number:  i + 1,
		}
		messageByte, err := json.MarshalIndent(message, "", "    ")
		if err != nil {
			log.Fatal("error while parsing message to json, err: " + err.Error())
		}
		err = ch.Publish(
			"",
			q.Name,
			false, false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        messageByte,
			},
		)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Published message ", i+1)
		time.Sleep(time.Second)
	}
	time.Sleep(time.Hour)
}
