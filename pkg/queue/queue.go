package queue

import (
	"context"
	"fmt"
	"log"
	"wt/pkg/enum"
	"wt/pkg/utils"

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

type TaskFailed struct {
	Id string
	TargerService string
	OriginatedBy string
	TaskName string
	DbUpdateValue string
}

type TaskSuccess struct {

}

type TaskStatus struct {
	Id string
	SentBy string
	TargetService string
	taskName string
	dbUpdateValue string
	taskStatus string
}

func NewTaskStatus(id string, sentBy string, targerService string, taskName string, dbUpdateValue string) *TaskStatus {
	return &TaskStatus{
		Id: id,
		SentBy: sentBy,
		TargetService: targerService,
		taskName: taskName,
		dbUpdateValue: dbUpdateValue,
	}
}

func (t *TaskStatus) ConvertToBytes() *[]byte {
	dataInBytes, _ := utils.ConvertIntoBytes(t)
	return dataInBytes
}

type ConsumerChan chan <-amqp.Delivery

func NewMessageQueue(conn *amqp.Connection, QueueName string) *MessageQueue {
	channel := createChannel(conn)
	queue := createQueue(channel, QueueName)

	return &MessageQueue{Ch: channel, queue: &queue}
}

func (m *MessageQueue) Publish(ctx context.Context, data *[]byte) error {
	msg := amqp.Publishing{
		ContentType: string(enum.ApplicationJsonType),
		Body:        *data,
		CorrelationId: string(enum.EmptyPlanCrrId),
	}

	err := m.Ch.PublishWithContext(ctx, "", m.queue.Name, false, false, msg)
	if err != nil {
		return fmt.Errorf("error in publishing : %w", err)
	}

	return nil
}


func (m *MessageQueue) Consume(queueName string) (<-chan amqp.Delivery, error) {
	consumerChan, err := m.Ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return consumerChan, fmt.Errorf("error occured while consuming from queue %v : %w", queueName, err)
	}

	return consumerChan, nil
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


