package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultPort = "8080"

func main() {

	// load any .env variables.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	// set the port.
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// set the queue URI.
	queueUri := os.Getenv("QUEUE_URI")
	if queueUri == "" {
		log.Printf("A queue URI is not available in the environment as \"QUEUE_URI\".\r\n")
		os.Exit(-1)
	}

	queueUser := os.Getenv("QUEUE_USERNAME")
	if queueUser == "" {
		log.Printf("A queue username is not available in the environment as \"QUEUE_USERNAME\".\r\n")
		os.Exit(-2)
	}

	queuePass := os.Getenv("QUEUE_PASSWORD")
	if queuePass == "" {
		log.Printf("A queue password is not available in the environment as \"QUEUE_PASSWORD\".\r\n")
		os.Exit(-3)
	}

	// connect to the broker.
	connection, err := amqp.Dial("amqp://" + queueUser + ":" + queuePass + "@" + queueUri)
	if err != nil {
		log.Printf("Failed to connect to the broker: %v\r\n", err.Error())
		os.Exit(-4)
	}
	defer connection.Close()

	// get a queue channel.
	channel, err := connection.Channel()
	if err != nil {
		log.Printf("Failed to get channel.\r\n")
		os.Exit(-5)
	}
	defer channel.Close()

	// declare the queue if necessary.
	queue, err := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Printf("Failed to create queue.\r\n")
		os.Exit(-6)
	}

	// define a context for the publishing.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to queue.\r\n")
		os.Exit(-7)
	}

	log.Printf(" [x] Sent %s\n", body)
}
