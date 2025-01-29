package rabbitMq

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"task-golang/model"
	"task-golang/service"
)

func InitRabbitMq(taskService *service.TaskService) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

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

	// Create a new logger instance
	logEntry := log.WithFields(log.Fields{
		"component": "RabbitMQ",
		"task":      "taskName",
	})

	// Create a context with logger and userId
	ctx := context.WithValue(context.Background(), model.ContextLogger, logEntry)
	ctx = context.WithValue(ctx, model.ContextUserID, int64(1)) // Set userId manually
	// Process messages in a separate goroutine
	go func() {
		for d := range msgs {
			var jsonString string
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
				log.Printf("Received a message: %+v", task)

				// Call CreateTask here
				errResponse := taskService.CreateTask(ctx, &task, 1)
				if errResponse != nil {
					log.Printf("Error creating task: %+v", errResponse)
				} else {
					log.Println("Task created successfully")
				}
			}
		}
	}()

	log.Println("RabbitMQ consumer started successfully and waiting for messages.")

	// **Keep the function running** by waiting indefinitely
	select {} // This prevents the function from exiting
}
