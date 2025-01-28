package config

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"task-golang/model"
)

func InitRabbitMq() {

	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare the queue
	q, err := ch.QueueDeclare(
		"q.nurlan.jpg", // queue name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	// Process messages
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var jsonString string
			// Unmarshal message body into a string
			err := json.Unmarshal(d.Body, &jsonString)
			if err != nil {
				log.Printf("Error decoding JSON string: %s", err)
				continue
			}

			// Now unmarshal the JSON string into the TaskRequestDto struct
			var task model.TaskRequestDto
			err = json.Unmarshal([]byte(jsonString), &task)
			if err != nil {
				log.Printf("Error decoding TaskRequestDto: %s", err)
			} else {
				fmt.Printf("Received a message: %+v\n", task)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
