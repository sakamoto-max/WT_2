package repository

// import (
// 	"context"
// 	// "errors"
// 	"exercise_service/internal/domain"
// 	"fmt"

// 	// "github.com/redis/go-redis/v9"
// 	// myerrors "github.com/sakamoto-max/wt_2-pkg/my_errors"
// )

// func (r *Repo) SetAllExercises(ctx context.Context, userId string, allExers *[]domain.Exercise) {
// 	// user_id:%v:all_exercises (main key)
// 	// exer_0 : exer_id | exercise_name | rest_time | body_part | equipment | created_at | updated_at
// 	// byExerId : exer_id exercise_name rest_time body_part equipment created_at updated_at
// 	// byExerName : exercise_name exer_id rest_time body_part equipment created_at updated_at

// 	mainKey := fmt.Sprintf("user_id:%s:all_exercises", userId)

// 	for _, eachExer := range *allExers {
// 		r.rDB.HSet(ctx, mainKey, "")
// 	}
// }

