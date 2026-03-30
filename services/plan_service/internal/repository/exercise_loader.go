package repository

// import (
// 	"context"
// 	"fmt"
// 	"plan_service/internal/models"

// 	"github.com/redis/go-redis/v9"
// )

// func (d *DBs) LoadExercises(ctx context.Context) error {

// 	allExercises, err := d.GetAllExercisesFromDB(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	err = d.SetAllExercises(ctx, allExercises)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (d *DBs) GetAllExercisesFromDB(ctx context.Context) (*[]models.Exercise, error) {

// 	var allExercises []models.Exercise

// 	rows, err := d.PDB.Query(ctx, `
// 		SELECT ID, EXERCISE_NAME FROM EXERCISES	
// 	`)
// 	if err != nil {
// 		return &allExercises, err
// 	}

// 	defer rows.Close()

// 	var id int
// 	var name string

// 	for rows.Next() {
// 		err := rows.Scan(&id, &name)
// 		if err != nil {
// 			return &allExercises, err
// 		}

// 		temp := models.Exercise{Id: id, Name: name}

// 		allExercises = append(allExercises, temp)
// 	}

// 	return &allExercises, nil
// }

// func (d *DBs) SetAllExercises(ctx context.Context, allExercises *[]models.Exercise) error {
// 	pipe := d.RDB.Pipeline()

// 	key1 := "exercise_name_to_exercise_id"
// 	key2 := "exercise_id_to_exericse_name"

// 	for _, v := range *allExercises {
// 		pipe.HSet(ctx, key1, v.Name, v.Id)
// 		pipe.HSet(ctx, key2, v.Id, v.Name)
// 	}

// 	_, err := pipe.Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (d *DBs) getAllUsersPlanNames(ctx context.Context, userId int) (*[]models.Plan3, error) {

// 	var planNamesWithIDs []models.Plan3

// 	rows, err := d.PDB.Query(ctx, `
// 		SELECT ID, PLAN_NAME FROM PLANS
// 		INNER JOIN USER_PLANS
// 		ON PLANS.ID = USER_PLANS.PLAN_ID
// 		WHERE USER_ID = $1
// 	`, userId)
// 	if err != nil {
// 		return &planNamesWithIDs, fmt.Errorf("error getting planNames from postgres : %w\n", err)
// 	}

// 	defer rows.Close()

// 	var id int
// 	var planName string

// 	for rows.Next() {

// 		err := rows.Scan(&id, &planName)
// 		if err != nil {
// 			return &planNamesWithIDs, err
// 		}

// 		a := models.Plan3{Id: id, PlanName: planName}

// 		planNamesWithIDs = append(planNamesWithIDs, a)
// 	}

// 	return &planNamesWithIDs, nil
// }

// func (d *DBs) SetAllUserPlanNames(ctx context.Context, userId int) error {
// 	key := fmt.Sprintf("user:%v:plan_names", userId)
// 	loadedKey := fmt.Sprintf("user:%v:plan_names_loaded", userId)

// 	planNamesWithIDs, err := d.getAllUsersPlanNames(ctx, userId)
// 	if err != nil {
// 		return err
// 	}

// 	pipe := d.RDB.Pipeline()

// 	for _, v := range *planNamesWithIDs {
// 		pipe.HSet(ctx, key, v.PlanName, v.Id)
// 	}

// 	pipe.Set(ctx, loadedKey, true, 0)

// 	_, err = pipe.Exec(ctx)
// 	if err != nil {
// 		return fmt.Errorf("error executing the pipe for setting all user plans : %w\n", err)
// 	}

// 	return nil
// }

// func (d *DBs) UserPlanNamesLoaded(ctx context.Context, userId int) (bool, error) {

// 	var exists bool

// 	existsKey := fmt.Sprintf("user:%v:plan_names_loaded", userId)

// 	err := d.RDB.Get(ctx, existsKey).Scan(&exists)
// 	if err != nil {
// 		if err == redis.Nil {
// 			return exists, nil
// 		}

// 		return exists, err
// 	}

// 	return exists, nil
// }

// func (d *DBs) ExerciseLoaded(ctx context.Context) error {
// 	// check if exercise_loaded is present
// 	// if not
// 	//    init exercise_loaded to false
// 	//    load the exercises
// 	// 	  set exercise_loaded to true
// 	// if present
// 	//    check if it is true or false
// 	//    if false
// 	//    load the exercises
// 	// 	  set exercise_loaded to true
// 	loaded, err := d.GetExerciseLoaded(ctx)
// 	if err != nil {
// 		return fmt.Errorf("error checking if exercises are loaded : %w\n", err)
// 	}

// 	if !loaded {
// 		err = d.LoadExercises(ctx)
// 		if err != nil {
// 			return fmt.Errorf("error loading the exercises to redis : %w\n", err)
// 		}
// 		err = d.SetExerciseLoaded(ctx)
// 		if err != nil {
// 			return fmt.Errorf("error setting exercise loaded : %w\n", err)
// 		}
// 	}
// 	return nil

// }
