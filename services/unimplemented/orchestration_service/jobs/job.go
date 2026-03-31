package jobs

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"orchestration_service/repository"
// 	"sync"
// 	"time"
// 	"wt/pkg/enum"
// 	"wt/pkg/utils"

// 	mq "wt/pkg/queue"

// 	"github.com/jackc/pgx/v5"
// )

// // create empty plan for the user after signing up

// func Operate(Db *repository.DB, planQueue *mq.MessageQueue, wg *sync.WaitGroup) {

// 	defer wg.Done()

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
// 	defer cancel()

// 	data, err := Db.FetchDataFromAuth(ctx)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			fmt.Println("no rows found")
// 			return
// 		}
// 	}

// 	dataInBytes, _ := utils.ConvertIntoBytes(data)

// 	err = planQueue.Publish(ctx, dataInBytes, string(enum.ApplicationJsonType))
// 	if err != nil {
// 		fmt.Printf("error occured while uploading data to the queue : %v", err)
// 		return
// 	}

// 	fmt.Printf("data sent")
// }
