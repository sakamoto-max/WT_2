package types

import "time"

type Data struct {
	Id            int       `db:"id"`
	TargetService string    `db:"target_service"`
	Task          string    `db:"task"`
	Status        string    `db:"status"`
	Payload       any       `db:"payload"`
	CreatedAt     time.Time `db:"created_at"`
}

