package queuemock

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueMock struct {
	Down    bool
	Channel chan amqp.Delivery
	Data    *[]byte
}

func (m *QueueMock) Publish(ctx context.Context, data *[]byte) error {
	if m.Down {
		return fmt.Errorf("queue is down")
	}

	fmt.Println("queue got data")

	m.Data = data

	return nil
}

func (m *QueueMock) Consume() (<-chan amqp.Delivery, error) {
	if m.Down {
		return nil, fmt.Errorf("queue is down")
	}

	return m.Channel, nil
}

func (m *QueueMock) Close() {

}
