package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
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
	wg := &sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go StartConsumer(i+1, wg)
	}
	wg.Wait()
}

func StartConsumer(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("New Consumer started id=%d\n", id)

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
		fmt.Sprintf("my consumer-%d", id),
		false,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	firstFalse := false
	for msg := range delivery {
		fmt.Printf("Customer=%d, message: %s\n", id, string(msg.Body))

		if id == 1 {
			time.Sleep(time.Second)
		} else {
			time.Sleep(time.Second * 4)
		}

		data := Data{}
		err := json.Unmarshal(msg.Body, &data)
		if err != nil {
			log.Fatal(err)
		}

		if !firstFalse && data.Number == 1 {
			firstFalse = true
			msg.Acknowledger.Reject(1, true)
			continue
		}
		
		err = msg.Acknowledger.Ack(msg.DeliveryTag, false)
		if err != nil {
			log.Fatal(err)
		}
	}
}
