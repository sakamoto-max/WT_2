package types

import (
	"time"
	"wt/pkg/utils"
)

type Data struct {
	Id            string    `db:"id"`
	TargetService string    `db:"target_service"`
	Task          string    `db:"task"`
	Status        string    `db:"status"`
	Payload       any       `db:"payload"`
	CreatedAt     time.Time `db:"created_at"`
	NumberOfTries int       `db:"number_of_tries"`
}

func (d *Data) Process() (*[]byte, error) {
	// convert into bytes
	dataInBytes, err := utils.ConvertIntoBytes(d)
	if err != nil {
		return nil, err
	}

	return dataInBytes, nil
}
