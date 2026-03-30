package main

// import (
// 	"log"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// func main() {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil{
// 		log.Fatalf("error creating consumer for plan : %w", err)
// 	}

// 	ch, err := conn.Channel()
// 	if err != nil{
// 		log.Fatalf("error creating channel : %w", err)
// 	}

// 	ch.QueueDeclare()

// }