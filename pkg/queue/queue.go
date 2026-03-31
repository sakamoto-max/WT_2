package queue

import (
	"context"
	"fmt"
	"log"
	"wt/pkg/enum"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConn() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("error opening a connection to rabbit mq : %v", err)
	}

	return conn
}

type MessageQueue struct {
	Ch    *amqp.Channel
	queue *amqp.Queue
}

type ConsumerChan chan <-amqp.Delivery

func NewMessageQueue(conn *amqp.Connection, QueueName string) *MessageQueue {
	channel := createChannel(conn)
	queue := createQueue(channel, QueueName)

	return &MessageQueue{Ch: channel, queue: &queue}
}

func (m *MessageQueue) Publish(ctx context.Context, data *[]byte, contentType string) error {

	msg := amqp.Publishing{
		ContentType: contentType,
		Body:        *data,
		CorrelationId: string(enum.EmptyPlanCrrId),
	}

	fmt.Println(msg)

	err := m.Ch.PublishWithContext(ctx, "", m.queue.Name, false, false, msg)
	if err != nil {
		return fmt.Errorf("error in publishing : %w", err)
	}

	return nil
}


func (m *MessageQueue) Consume(queueName string, chanName <- chan amqp.Delivery) <-chan amqp.Delivery {
	chanName, err := m.Ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("error occured : %v", err)
	}

	// chanName = msgs



	return chanName
}



func createChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error in creating a channel : %v", err)
	}

	return ch
}
func createQueue(ch *amqp.Channel, queueName string) amqp.Queue {
	queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error creating %v : %v", queue.Name, err)
	}
	return queue
}


