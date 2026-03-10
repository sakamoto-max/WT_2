package repository

import (
	"context"
	"fmt"
	customerrors "plan_service/internal/custom_errors"
	"strconv"

	"github.com/redis/go-redis/v9"
)

//
// KEYS :
// workout ongoing := fmt.Sprintf("user:%v:workout_ongoing", userId)
// current workout Id := fmt.Sprintf("user:%v:workout_id", userId)
// New Plan Name := fmt.Sprintf("user:%v:new_plan_name", userId)
// setting exercises to plan Name := fmt.Sprintf("user:%v:%v:%v", userId, planName, i(number of the exercise))
// exercise_id_list := fmt.Sprintf("user:%v:exercise_id_list", userId)
// current exercise id := fmt.Sprintf("user:%v:current_exercise_id", userID)
// current set := fmt.Sprintf("user:%v:current_set", userId)
// exercise loaded := "exercise_loaded"
// all exercises := "all exercises"
// number of sets := fmt.Sprintf("user:%v:exercise_id:%v:no_of_sets", userId, exerciseId)
// for hset reps and weight :=  fmt.Sprintf("user:%v:exercise_id:%v", userId, exerciseId)
// rep field := fmt.Sprintf("set:%v:reps", setNumber)
// weight field := fmt.Sprintf("set:%v:weight", setNumber)

func (d *DBs) OngoinWorkoutExists(ctx context.Context, userId int) (bool, error) {

	var exists string

	key := fmt.Sprintf("user:%v:workout_ongoing", userId)

	err := d.RDB.Get(ctx, key).Scan(&exists)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}
	return true, nil
}
func (d *DBs) DelWorkoutOngoing(ctx context.Context, userID int) error {
	workoutOngoingKey := fmt.Sprintf("user:%v:workout_ongoing", userID)

	err := d.RDB.Del(ctx, workoutOngoingKey).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) SetWorkoutId(ctx context.Context, userId int, workoutID int) error {
	key := fmt.Sprintf("user:%v:workout_id", userId)

	err := d.RDB.Set(ctx, key, workoutID, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) SetOngoingWorkout(ctx context.Context, userId int) error {
	key := fmt.Sprintf("user:%v:workout_ongoing", userId)
	err := d.RDB.Set(ctx, key, true, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) SetNewPlanName(ctx context.Context, userId int, planName string) error {
	key := fmt.Sprintf("user:%v:new_plan_name", userId)

	err := d.RDB.Set(ctx, key, planName, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) AppendToExerciseList(ctx context.Context, userId int, exericeId int) error {
	key := fmt.Sprintf("user:%v:exercise_id_list", userId)

	err := d.RDB.LPush(ctx, key, exericeId).Err()
	if err != nil {
		return fmt.Errorf("error pushing exercise id into exercise list : %w", err)
	}

	return nil
}
func (d *DBs) GetPlanName(ctx context.Context, userID int) (string, error) {
	key := fmt.Sprintf("user:%v:new_plan_name", userID)

	var planName string
	err := d.RDB.Get(ctx, key).Scan(&planName)
	if err != nil {
		return planName, err
	}
	return planName, nil
}
func (d *DBs) SetExercise(ctx context.Context, userId int, planName string, exerciseNames *[]string) error {
	pipe := d.RDB.Pipeline()

	for i, v := range *exerciseNames {

		key := fmt.Sprintf("user:%v:%v:%v", userId, planName, i)

		pipe.Set(ctx, key, v, 0)

	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) GetExerciseLoaded(ctx context.Context) (bool, error) {
	key := "exercise_loaded"
	var val bool

	err := d.RDB.Get(ctx, key).Scan(&val)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return val, err
	}

	return val, nil
}
func (d *DBs) SetExerciseLoaded(ctx context.Context) error {
	key := "exercise_loaded"
	err := d.RDB.Set(ctx, key, true, 0).Err()

	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) DelExerciseLoaded(ctx context.Context) error {
	key := "exercise_loaded"
	err := d.RDB.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) GetExerciseIdFromMain(ctx context.Context, exerciseName string) (int, error) {
	var id int

	key := "exercise_name_to_exercise_id"

	err := d.RDB.HGet(ctx, key, exerciseName).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}
func (d *DBs) ExerciseExistsInMain(ctx context.Context, exerciseName string) (bool, error) {
	var id int

	key := "exercise_name_to_exercise_id"

	err := d.RDB.HGet(ctx, key, exerciseName).Scan(&id)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, fmt.Errorf("error checking if the exercise %v exists in main : %w\n", exerciseName, err)
	}

	return true, err
}
func (d *DBs) GetWorkoutId(ctx context.Context, userID int) (int, error) {
	key := fmt.Sprintf("user:%v:workout_id", userID)

	var id int
	err := d.RDB.Get(ctx, key).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("error getting workout tracker id from redis : %w\n", err)
	}

	return id, nil
}
func (d *DBs) DelWorkoutId(ctx context.Context, userId int) error {
	key := fmt.Sprintf("user:%v:workout_id", userId)

	err := d.RDB.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) SetCurrentExerciseID(ctx context.Context, userId int, exerciseId int) error {
	key := fmt.Sprintf("user:%v:current_exercise_id", userId)

	err := d.RDB.Set(ctx, key, exerciseId, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) GetCurrentExerciseId(ctx context.Context, userID int) (int, error) {
	var id int

	key := fmt.Sprintf("user:%v:current_exercise_id", userID)

	err := d.RDB.Get(ctx, key).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("error getting current exercise_id from redis : %w\n", err)
	}

	return id, nil
}
func (d *DBs) DelCurrentExerciseId(ctx context.Context, userId int) error {
	key := fmt.Sprintf("user:%v:current_exercise_id", userId)

	err := d.RDB.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting current_exercise_id : %w\n", err)
	}

	return nil

}
func (d *DBs) SetCurrentSet(ctx context.Context, userId int) error {
	key := fmt.Sprintf("user:%v:current_set", userId)

	err := d.RDB.Set(ctx, key, 1, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) IncrCurrentSet(ctx context.Context, userId int) error {

	key := fmt.Sprintf("user:%v:current_set", userId)

	err := d.RDB.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) GetCurrentSet(ctx context.Context, userID int) (int, error) {

	var currentSet int

	key := fmt.Sprintf("user:%v:current_set", userID)

	err := d.RDB.Get(ctx, key).Scan(&currentSet)
	if err != nil {
		return currentSet, err
	}

	return currentSet, nil
}
func (d *DBs) DelCurrentSet(ctx context.Context, userId int) error {

	key := fmt.Sprintf("user:%v:current_set", userId)

	err := d.RDB.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
func (d *DBs) CurrentSetExists(ctx context.Context, userId int) (bool, error) {

	var currentSet int

	key := fmt.Sprintf("user:%v:current_set", userId)

	err := d.RDB.Get(ctx, key).Scan(&currentSet)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (d *DBs) SetRepsWeight(ctx context.Context, userId int, exerciseId int, setNumber int, reps int, weight int) error {
	var noOfSetsDone int

	noOfSets := fmt.Sprintf("user:%v:exercise_id:%v:no_of_sets", userId, exerciseId)

	err := d.RDB.Get(ctx, noOfSets).Scan(&noOfSetsDone)
	if err != nil {
		return fmt.Errorf("error getting number of sets done for %v : %w\n", exerciseId, err)
	}

	// reps :
	// user:id:exercise_id:id hset set:1:reps reps

	// weight :
	// user:id:exercise_id:id hset set:1:weight weight

	// no of sets :
	// incr
	// user:id:exercise_id:id:no_of_sets 1

	key := fmt.Sprintf("user:%v:exercise_id:%v", userId, exerciseId)
	repField := fmt.Sprintf("set:%v:reps", setNumber)
	weightsField := fmt.Sprintf("set:%v:weight", setNumber)

	pipe := d.RDB.Pipeline()

	pipe.HSet(ctx, key, repField, reps, weightsField, weight)
	pipe.Set(ctx, noOfSets, noOfSetsDone+1, 0)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error setting reps and weight : %w\n", err)
	}

	return nil

}
func (d *DBs) SetNoOfSets(ctx context.Context, userId int, exerciseId int) error {

	key := fmt.Sprintf("user:%v:exercise_id:%v:no_of_sets", userId, exerciseId)

	err := d.RDB.Set(ctx, key, 0, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting no of sets : %w\n", err)
	}

	return nil
}

func (d *DBs) IncrNoOfSets(ctx context.Context, userId int, exerciseId int) error {

	key := fmt.Sprintf("user:%v:exercise_id:%v:no_of_sets", userId, exerciseId)

	err := d.RDB.Incr(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error occured while increasing no_of_sets : %w", err)
	}

	return nil
}

func (d *DBs) GetNumberOfSets(ctx context.Context, userId int, exerciseId int) (int, error) {

	var noOfSetsDone int

	key := fmt.Sprintf("user:%v:exercise_id:%v:no_of_sets", userId, exerciseId)

	err := d.RDB.Get(ctx, key).Scan(&noOfSetsDone)

	if err != nil {
		return noOfSetsDone, fmt.Errorf("error getting number of sets done from redis : %v\n", err)
	}
	return noOfSetsDone, nil
}

func (d *DBs) GetWeightForSet(ctx context.Context, userId int, exerciseId int, setNumber int) (int, error) {

	var weight int

	key := fmt.Sprintf("user:%v:exercise_id:%v", userId, exerciseId)
	weightsField := fmt.Sprintf("set:%v:weight", setNumber)

	err := d.RDB.HGet(ctx, key, weightsField).Scan(&weight)
	if err != nil {
		return weight, fmt.Errorf("error getting weight from redis : %v\n", err)
	}

	return weight, nil
}
func (d *DBs) GetRepsForSet(ctx context.Context, userId int, exerciseId int, setNumber int) (int, error) {

	var reps int

	key := fmt.Sprintf("user:%v:exercise_id:%v", userId, exerciseId)
	repField := fmt.Sprintf("set:%v:reps", setNumber)

	err := d.RDB.HGet(ctx, key, repField).Scan(&reps)
	if err != nil {
		return reps, fmt.Errorf("error getting the reps from redis : %v\n", err)
	}

	return reps, nil
}

func (d *DBs) ResetCurrentSet(ctx context.Context, userId int) error {
	key := fmt.Sprintf("user:%v:current_set", userId)

	err := d.RDB.Set(ctx, key, 1, 0).Err()
	if err != nil {
		return fmt.Errorf("error resetting the current set : %v\n", err)
	}

	return nil
}

func (d *DBs) EndWorkoutred2(ctx context.Context, userId int) ([]string, error) {

	exerciseList := fmt.Sprintf("user:%v:exercise_id_list", userId)

	idList, err := d.RDB.LRange(ctx, exerciseList, 0, -1).Result()
	if err != nil {
		//
		return idList, err
	}
	return idList, nil

}

func (d *DBs) GetExerciseIDList(ctx context.Context, userId int) ([]int, error) {

	var exerIdList []int

	exerciseList := fmt.Sprintf("user:%v:exercise_id_list", userId)

	idList, err := d.RDB.LRange(ctx, exerciseList, 0, -1).Result()
	if err != nil {
		//
		return exerIdList, fmt.Errorf("error getting exerciseId List : %w", err)
	}

	for _, v := range idList {
		a, err := strconv.Atoi(v)
		if err != nil {
			return exerIdList, fmt.Errorf("error converting string to int while getting exercise Id list %w", err)
		}

		exerIdList = append(exerIdList, a)
	}

	return exerIdList, nil
}

func (d *DBs) GetNoOfSetsForExer(ctx context.Context, userId int, exerciseId int) (int, error) {

	var totalSets int

	noOfSetsKey := fmt.Sprintf("user:%v:exercise_id:%v:no_of_sets", userId, exerciseId)

	noOfSets, err := d.RDB.Get(ctx, noOfSetsKey).Result()
	if err != nil {
		return totalSets, fmt.Errorf("error while getting no of sets for exercise id : %v : %w/n", exerciseId, err)
	}

	totalSets, err = strconv.Atoi(noOfSets)
	if err != nil {
		return totalSets, fmt.Errorf("error while converting no of sets from string to int while getting total sets : %w", err)
	}

	return totalSets, nil
}

func (d *DBs) DelUserInfo(ctx context.Context, userId int, exerciseIDList *[]int) error {

	pipe := d.RDB.Pipeline()

	for _, v := range *exerciseIDList {
		setsRepsKey := fmt.Sprintf("user:%v:exercise_id:%v", userId, v)
		noOfSets := fmt.Sprintf("user:%v:exercise_id:%v:no_of_sets", userId, v)

		pipe.Del(ctx, setsRepsKey)
		pipe.Del(ctx, noOfSets)
	}

	exerciseIDListKey := fmt.Sprintf("user:%v:exercise_id_list", userId)
	workoutOngoing := fmt.Sprintf("user:%v:workout_ongoing", userId)
	currentExerIdKey := fmt.Sprintf("user:%v:current_exercise_id", userId)
	currentSetKey := fmt.Sprintf("user:%v:current_set", userId)
	workoutId := fmt.Sprintf("user:%v:workout_id", userId)

	pipe.Del(ctx, exerciseIDListKey)
	pipe.Del(ctx, workoutOngoing)
	pipe.Del(ctx, currentExerIdKey)
	pipe.Del(ctx, currentSetKey)
	pipe.Del(ctx, workoutId)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting user Info : %w", err)
	}

	return nil
}

func (d *DBs) GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error) {

	key := "exercise_id_to_exericse_name"

	exerciserName, err := d.RDB.HGet(ctx, key, exerciseId).Result()
	if err != nil {
		return exerciserName, err
	}

	return exerciserName, nil
}

func (d *DBs) GetPlanIdFromRedis(ctx context.Context, userId int, planName string) (int, error) {

	var id int

	key := fmt.Sprintf("user:%v:plan_names", userId)

	err := d.RDB.HGet(ctx, key, planName).Scan(&id)

	if err != nil {
		if err == redis.Nil {
			return id, customerrors.ErrPlanNameDoesNotExists
		}

		return id, err
	}

	return id, nil
}

func (d *DBs) SetWorkoutWithPlan(ctx context.Context, userId int) error {

	key := fmt.Sprintf("user:%v:workout_with_plan", userId)

	err := d.RDB.Set(ctx, key, true, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// func SetExercisesForPlan(userID int, planName string ,exerciseNames *[]string) {

// 	key := fmt.Sprintf("user:%v:plan_name%v:exercises", userID, planName)

// }
