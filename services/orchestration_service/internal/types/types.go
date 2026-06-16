package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Data struct {
	TaskId        string
	WorkersTries  int
	TimeOut       timeOut
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

func (d *Data) InitTimeOut() {
	d.TimeOut = timeOut{
		Level:    1,
		Duration: time.Minute * 1,
	}
}

func (d *Data) IncreaseTimeOut() {
	d.TimeOut = timeOut{
		Level:    d.TimeOut.Level + 1,
		Duration: d.TimeOut.Duration * 2,
	}
}



type timeOut struct {
	Level    int
	Duration time.Duration
}

// task Id -> gen random UUID, put it in context and send
// if push failes -> send it back to the queue with a default timeout
// when the worker receives the data check if it has timeout set in redis -> if not try to push it
// if it has timer set -> skip it

// if the push failes again -> double the timer
// if the push failes after 5 tries -> task failed
