package mock

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
)

type MockQueue struct {
	Open bool
}

func (m *MockQueue) Publish(ctx context.Context, data *[]byte) error {
	if !m.Open {
		return fmt.Errorf("queue is closed")
	}

	return nil
}

func (m *MockQueue) Consume() (<-chan amqp.Delivery, error) {

	if !m.Open {
		return nil, fmt.Errorf("queue is closed")
	}

	c := make(chan amqp.Delivery)

	data := mqTypes.Data{
		DbId:          "123",
		TaskName:      "mock task",
		SentBy:        "mock service",
		TaskStatus:    "mock status",
		TargetService: "mock target",
	}

	dataInBytes, _ := data.ConvertIntoBytes()

	d := amqp.Delivery{
		Body: *dataInBytes,
	}

	c <- d
	c <- d
	c <- d

	return c, nil
}
