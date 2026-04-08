module email_service

go 1.25.4

require (
	go.uber.org/zap v1.27.1
	wt/pkg v0.0.0
)

require (
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

replace wt/pkg => ../../pkg
