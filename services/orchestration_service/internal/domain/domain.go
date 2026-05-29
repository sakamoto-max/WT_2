package domain

import (
	"fmt"
	"orchestration_service/internal/utils"
	"time"
)


type Data struct {
	DbId          string    `db:"id"`
	TargetService string    `db:"target_service"`
	CreatedBy     string    `db:"created_by"`
	Task          string    `db:"task"`
	Status        string    `db:"status"`
	Payload       any       `db:"payload"`
	CreatedAt     time.Time `db:"created_at"`
	NumberOfTries int       `db:"number_of_tries"`
}

func (d *Data) ConvertToBytes() (*[]byte, error) {

	dataInBytes, err := utils.ConvertIntoBytes(d)
	if err != nil {
		return nil, fmt.Errorf("error while converting data into bytes : %w", err)
	}

	return dataInBytes, nil
}
