package types

import (
	"fmt"
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
	NumberOfTries *int      `db:"number_of_tries"`
}

func (d *Data) ConvertToBytes() (*[]byte, error) {
	// convert into bytes
	dataInBytes, err := utils.ConvertIntoBytes(d)
	if err != nil {
		return nil, fmt.Errorf("error while converting data into bytes : %w", err)
	}

	return dataInBytes, nil
}

