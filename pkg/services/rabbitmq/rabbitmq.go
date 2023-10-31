package rabbitmq

import (
	"fmt"
	"github.com/esirangelomub/go-chat-application/configs"
	entityPkg "github.com/esirangelomub/go-chat-application/pkg/models/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func SetupRabbitMQ(config *configs.Conf) (*amqp.Connection, *amqp.Channel) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		config.RabbitMQUser,
		config.RabbitMQPassword,
		config.RabbitMQHost,
		config.RabbitMQPort,
	)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	qb, err := ch.QueueDeclare(
		config.RabbitMQQueueBot,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue Bot: %v", err)
	}

	qw, err := ch.QueueDeclare(
		config.RabbitMQQueueWebSocket,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue WebSocket: %v", err)
	}

	// Declare the Exchange
	err = ch.ExchangeDeclare(
		config.RabbitMQExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	// Bind the queues to the exchange
	err = ch.QueueBind(
		qb.Name,
		"bot",
		config.RabbitMQExchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue Bot: %v", err)
	}

	err = ch.QueueBind(
		qw.Name,
		"websocket",
		config.RabbitMQExchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue WebSocket: %v", err)
	}

	return conn, ch
}

func Consume(ch *amqp.Channel, out chan<- amqp.Delivery, queue string) error {
	log.Printf("Starting consumer for queue: %s", queue)

	msgs, err := ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("Consumer registered for queue: %s", queue)

	for msg := range msgs {
		log.Printf("Received a message from queue: %s", queue)
		out <- msg
	}
	
	return nil
}

func Publish(ch *amqp.Channel, msg entityPkg.ChatMessage, exName, key string) error {
	body, err := msg.ToJSON()
	if err != nil {
		return err
	}
	err = ch.Publish(
		exName,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
