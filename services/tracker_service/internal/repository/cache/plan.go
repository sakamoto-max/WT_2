package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type planCache struct {
	client *redis.Client
}

func (c *planCache) SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error {

	key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	noOfExercisesKey := "number_of_exercises"

	pipe := c.client.Pipeline()

	for i, exerciseName := range *exerciseNames {
		field := fmt.Sprintf("exer_%v", i)
		pipe.HSet(ctx, key, field, exerciseName)
	}

	pipe.HSet(ctx, key, noOfExercisesKey, len(*exerciseNames))

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error in setting plan in redis : %w", err)
	}

	return nil
}
func (c *planCache) GetPlanWithExercises(ctx context.Context, userId string, planName string) (*[]string, error) {

	key := fmt.Sprintf("user_id:%v:plan_name:%v", userId, planName)
	noOfExercisesKey := "number_of_exercises"

	var NumberOfExercises int
	var numberOfExercisesInString string

	err := c.client.HGet(ctx, key, noOfExercisesKey).Scan(&numberOfExercisesInString)
	if err != nil {
		return nil, fmt.Errorf("error getting the number of exercises from redis : %v", err)
	}

	NumberOfExercises, err = strconv.Atoi(numberOfExercisesInString)
	if err != nil {
		return nil, fmt.Errorf("error converting to integer : %w", err)
	}

	var allExercises []string
	var exerciseName string

	for i := 0; i < NumberOfExercises; i++ {
		exerfield := fmt.Sprintf("exer_%v", i)
		err := c.client.HGet(ctx, key, exerfield).Scan(&exerciseName)
		if err != nil {
			return nil, fmt.Errorf("unable to get the exercise name : %w", err)
		}

		allExercises = append(allExercises, exerciseName)
	}

	return &allExercises, nil
}
