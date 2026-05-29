package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type newExerCache struct {
	client *redis.Client
}

func (c *newExerCache) SetUserNewExercises(ctx context.Context, userId string, exerciseNames *[]string) error {

	key := fmt.Sprintf("user_id:%v:new_exercises", userId)
	noOfExercisesKey := "number_of_exercises"

	pipe := c.client.Pipeline()

	pipe.HSet(ctx, key, noOfExercisesKey, len(*exerciseNames))

	for i, exerciseName := range *exerciseNames {
		exerciseKey := fmt.Sprintf("exercise_%v", i)
		pipe.HSet(ctx, key, exerciseKey, exerciseName)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (c *newExerCache) GetUserNewExercises(ctx context.Context, userId string) (*[]string, error) {

	key := fmt.Sprintf("user_id:%v:new_exercises", userId)
	noOfExercisesKey := "number_of_exercises"

	var noOfExericisesString string

	err := c.client.HGet(ctx, key, noOfExercisesKey).Scan(&noOfExericisesString)
	if err != nil {
		return nil, fmt.Errorf("error getting number of exercises : %w", err)
	}

	noOfExerises, err := strconv.Atoi(noOfExericisesString)
	if err != nil {
		return nil, fmt.Errorf("error converting no of exercises from string to integer : %w", err)
	}

	var allExercises []string

	for i := range noOfExerises {

		var oneExercise string

		exerciseKey := fmt.Sprintf("exercise_%v", i)
		err := c.client.HGet(ctx, key, exerciseKey).Scan(&oneExercise)
		if err != nil {
			return nil, fmt.Errorf("error getting the exercise name : %w", err)
		}

		allExercises = append(allExercises, oneExercise)

	}

	return &allExercises, nil
}
