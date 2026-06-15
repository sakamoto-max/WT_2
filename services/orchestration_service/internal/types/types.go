package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Data struct {
	DbId          string         `db:"id"`
	TargetService string         `db:"target_service"`
	CreatedBy     string         `db:"created_by"`
	Task          string         `db:"task"`
	Status        string         `db:"status"`
	Payload       map[string]any `db:"payload"`
	CreatedAt     time.Time      `db:"created_at"`
	NumberOfTries int            `db:"number_of_tries"`
	ServiceName   string
	NoData        bool
	Err           error
}

func (d *Data) ConvertToBytes() (*[]byte, error) {

	dataInBytes, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("error in converting data into bytes : %w", err)
	}

	return &dataInBytes, nil
}
