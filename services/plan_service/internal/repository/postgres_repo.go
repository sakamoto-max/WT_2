package repository

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/models"

	// "plan_service/internal/models"

	"github.com/jackc/pgx/v5"
)

func (d *DBs) PlanExists(ctx context.Context, userId int, planName string) (bool, error) {
	var planId int
	err := d.PDB.QueryRow(ctx, `
		select id from plans
		WHERE user_id = $1 AND NAME = $2
	`, userId, planName).Scan(&planId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if the plan already exists : %w\n", err)
	}

	return true, nil
}
func (d *DBs) PlanExistsReturnsId(ctx context.Context, userId int, planName string) (bool, int, error) {
	var planId int
	err := d.PDB.QueryRow(ctx, `
		select id from plans
		WHERE user_id = $1 AND NAME = $2
	`, userId, planName).Scan(&planId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, planId, nil
		}
		return false, planId, fmt.Errorf("error checking if the plan already exists : %w\n", err)
	}

	return true, planId, nil
}
func (d *DBs) CreatePlan(ctx context.Context, userId int, planName string, exerciseIds []int) error {

	var planId int

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaction : %w\n", err)
	}

	defer trnx.Rollback(ctx)

	err = trnx.QueryRow(ctx, `	
		INSERT INTO PLANS(USER_ID, NAME, CREATED_AT)
		VALUES($1, $2, NOW())	
		RETURNING ID
	`, userId, planName).Scan(&planId)
	if err != nil {
		return fmt.Errorf("error inserting plan name into plans : %w\n", err)
	}

	for _, exerciseId := range exerciseIds {
		_, err := trnx.Exec(ctx, `
			INSERT INTO PLAN_EXERCISES(PLAN_ID, EXERCISE_ID)
			VALUES($1, $2)		
		`, planId, exerciseId)
		if err != nil {
			return fmt.Errorf("error inserting exercise_id %v into plan_exercises : %w", exerciseId, err)
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting the trnx for creating plan : %w", err)
	}

	return nil
}
func (d *DBs) GetAllExercisesByPlanID(ctx context.Context, planId int) (*[]int, error) {
	var exerciseIDs []int

	rows, err := d.PDB.Query(ctx, `
		SELECT EXERCISE_ID FROM PLAN_EXERCISES
		WHERE PLAN_ID = $1
	`, planId)
	if err != nil {
		return &exerciseIDs, err
	}

	defer rows.Close()

	var id int

	for rows.Next() {

		err := rows.Scan(&id)
		if err != nil {
			return &exerciseIDs, err
		}

		exerciseIDs = append(exerciseIDs, id)
	}

	return &exerciseIDs, nil
}
func (d *DBs) AddExercisesToPlan(ctx context.Context, planId int, exerciseIDs *[]int) error {
	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	for _, id := range *exerciseIDs {
		_, err := trnx.Exec(ctx, `
			INSERT INTO plan_exercises(plan_id, exercise_id)
			VALUES ($1, $2)
		`, planId, id)
		if err != nil {
			return fmt.Errorf("error inserting the exericse with id %v : %w", id, err)
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing : %w", err)
	}

	return nil
}
func (d *DBs) DeleteExerciseFromPlan(ctx context.Context, planId int, exerciseIDs *[]int) error {

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	for _, id := range *exerciseIDs {
		_, err := trnx.Exec(ctx, `
			DELETE FROM plan_exercises
			WHERE PLAN_ID = $1 AND EXERCISE_ID = $2
		`, planId, id)
		if err != nil {
			return fmt.Errorf("error deleting exercise with id %v : %w", id, err)
		}

	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing : %w", err)
	}

	return nil
}
func (d *DBs) GetAllUserPlans(ctx context.Context, userId int) (*[]models.Plan3, error) {

	var allPlans []models.Plan3

	rows, err := d.PDB.Query(ctx, `
		SELECT ID, NAME FROM PLANS
		WHERE USER_ID = $1
	`, userId)

	if err != nil {
		return &allPlans, err
	}

	var id int
	var planName string

	for rows.Next() {

		err := rows.Scan(&id, &planName)
		if err != nil {
			return &allPlans, fmt.Errorf("error scaning rows : %w", err)
		}

		a := models.Plan3{PlanName: planName, Id: id}

		allPlans = append(allPlans, a)
	}

	rows.Close()

	return &allPlans, nil
}
func (d *DBs) DeletePlan(ctx context.Context, userId int, planId int) error {
	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaction : %w", err)
	}

	_, err = trnx.Exec(ctx, `
		DELETE FROM plan_exercises
		WHERE PLAN_ID = $1
	`, planId)
	if err != nil {
		return fmt.Errorf("error deleting exercises : %w", err)
	}

	_, err = trnx.Exec(ctx, `
		DELETE FROM PLANS
		WHERE ID = $1	
	`, planId)
	if err != nil {
		return fmt.Errorf("error deleting plan : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error committing : %w", err)
	}

	return nil

}

func (d *DBs) GetEmptyPlanID(ctx context.Context, userId int) (int, error) {
	var emptyPlanId int
	err := d.PDB.QueryRow(ctx, `
		SELECT id FROM plans
		WHERE name = 'empty' AND user_id = $1;
	`, userId).Scan(&emptyPlanId)
	if err != nil{
		return emptyPlanId, fmt.Errorf("error getting empty plan Id : %w", err)
	}

	return emptyPlanId, nil

}

func (d *DBs) PlanExistsReturnPlan(ctx context.Context, userID int, planName string) (bool, int, *[]int, error) {

	var allExerciseIDs *[]int
	exists, planId, err :=  d.PlanExistsReturnsId(ctx, userID, planName)
		if err != nil {
			return exists, planId, allExerciseIDs, fmt.Errorf("error getting exercise ids : %w", err)
	}
	if !exists {
		return exists, planId, allExerciseIDs, nil
	}

	allExerciseIDs, err = d.GetAllExercisesByPlanID(ctx, planId)
	if err != nil {
		return exists, planId, allExerciseIDs, fmt.Errorf("error getting exercise ids : %w", err)
	}

	return exists, planId, allExerciseIDs, nil
}

func (d *DBs) CreateEmptyPlan(ctx context.Context, userId int) (error) {
	_, err := d.PDB.Exec(ctx, `
		INSERT INTO plans(user_id, name, created_at)
		VALUES($1, $2, NOW())
	`, userId, models.EmptyPlan)
	if err != nil{
		return fmt.Errorf("Error creating empty plan : %v", err)
	}

	return nil
}







func (d *DBs) EmptyPlanExists(ctx context.Context, userId int) (bool, error) {
	var id int

	err := d.PDB.QueryRow(ctx, `
		SELECT PLANS.ID FROM USER_PLANS
		INNER JOIN PLANS
		ON USER_PLANS.PLAN_ID = PLANS.ID
		WHERE PLAN_NAME = 'empty' AND USER_ID = $1		
	`, userId).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *DBs) AddPlanAndExercisesToDB(ctx context.Context, userId int, planName string, exerciseNames *[]string) error {

	var planId int

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return err
	}

	defer trnx.Rollback(ctx)

	// add planName to plans
	err = trnx.QueryRow(ctx, `
		INSERT INTO PLANS(PLAN_NAME, CREATED_AT)
		VALUES($1, NOW())
		RETURNING ID	
	`, planName).Scan(&planId)
	if err != nil {
		return err
	}
	// add plan_id, user_id to user_plans

	_, err = trnx.Exec(ctx, `
		INSERT INTO USER_PLANS(USER_ID, PLAN_ID)
		VALUES($1, $2)	
	`, userId, planId)
	if err != nil {
		return err
	}

	err = d.addEachExercise(ctx, trnx, planId, exerciseNames)
	if err != nil {
		return err
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil

}

func (d *DBs) addEachExercise(ctx context.Context, trnx pgx.Tx, planId int, exerciseNames *[]string) error {
	for _, v := range *exerciseNames {
		// get the id of the exercise
		id, err := d.GetExerciseIdFromMain(ctx, v)
		if err != nil {
			return err
		}

		// insert plan_id, exercise_id and rest_time
		_, err = trnx.Exec(context.Background(), `
			INSERT INTO PLAN_EXERCISES(PLAN_ID, EXERCISE_ID, REST_TIME_IN_SECONDS)
			VALUES($1, $2, 120)		
		`, planId, id)

		if err != nil {
			return err
		}

	}
	return nil
}


func (d *DBs) GetPlanId(ctx context.Context, userId int, planName string) (int, error) {
	var planId int

	err := d.PDB.QueryRow(ctx, `
		SELECT PLANS.ID FROM USER_PLANS
		INNER JOIN PLANS
		ON USER_PLANS.PLAN_ID = PLANS.ID
		WHERE PLAN_NAME = $1 AND USER_ID = $2
	`, planName, userId).Scan(&planId)

	if err != nil {
		return planId, err
	}

	return planId, nil
}

func (d *DBs) GetPlanNameByID(ctx context.Context, planId int) (string, error) {
	var planName string

	err := d.PDB.QueryRow(ctx, `
		SELECT PLAN_NAME FROM PLANS
		WHERE ID = $1	
	`, planId).Scan(&planName)

	if err != nil {
		return planName, err
	}

	return planName, nil
}

// func (d *DBs) GetEmptyPlanID(ctx context.Context, userId int) (int, error) {
// 	var planId int

// 	err := d.PDB.QueryRow(ctx, `
// 		SELECT PLANS.ID FROM USER_PLANS
// 		INNER JOIN PLANS
// 		ON USER_PLANS.PLAN_ID = PLANS.ID
// 		WHERE PLAN_NAME = 'empty' AND USER_ID = $1
// 	`, userId).Scan(&planId)

// 	if err != nil {
// 		return planId, err
// 	}

// 	return planId, nil
// }

func (d *DBs) StartNewWorkout(ctx context.Context, userId int, planId int) (int, error) {
	var id int

	err := d.PDB.QueryRow(ctx, `
		INSERT INTO WORKOUT_TRACKER(USER_ID, PLAN_ID, STARTED_AT)
		VALUES($1, $2, NOW())	
		RETURNING ID
	`, userId, planId).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil

}

func (d *DBs) EndWorkoutPost(ctx context.Context, workoutId int) error {
	_, err := d.PDB.Exec(ctx, `
	 	UPDATE WORKOUT_TRACKER
		SET ENDED_AT = NOW()
		WHERE ID = $1
	`, workoutId)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBs) PushRepWeights(ctx context.Context, userId int, exerciseIdList *[]int, queries *[]string) error {
	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error begining the transactions : %w", err)
	}

	defer trnx.Rollback(ctx)

	for _, v := range *queries {
		_, err := trnx.Exec(ctx, v)
		if err != nil {
			return fmt.Errorf("error commiting in postgres : %w", err)
		}
	}

	err = d.DelUserInfo(ctx, userId, exerciseIdList)
	if err != nil {
		trnx.Rollback(ctx)
		return fmt.Errorf("error commiting in redis : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting : %w", err)
	}

	return nil
}

func (d *DBs) NoOfWorkoutsExistsInP() {

}
