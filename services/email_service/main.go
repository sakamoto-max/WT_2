package main

import (
	"log"
	"wt/pkg/enum"
	"wt/pkg/queue"
	"wt/pkg/utils"
)

func main() {

	conn := queue.NewConn()

	emailQueue := queue.NewMessageQueue(conn, string(enum.EmailQueue))

	msgs, err := emailQueue.Consume(string(enum.EmailQueue))
	if err != nil {
		log.Fatalf("error consuming from the email queue : %v", err)
	}

	log.Printf("email consumer has started")

	for msg := range msgs {
		data := utils.ConvertIntoJosn(&msg.Body)
		email := data.Payload["email"]
		log.Printf("sending email to : %v", email)
		continue
	}
}