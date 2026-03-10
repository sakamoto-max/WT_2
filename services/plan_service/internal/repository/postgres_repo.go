package repository

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/models"

	"github.com/jackc/pgx/v5"
)

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

func (d *DBs) CreateEmptyPlan(ctx context.Context, userId int) error {
	var planId int

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return err
	}

	defer trnx.Commit(ctx)

	err = trnx.QueryRow(ctx, `
		INSERT INTO PLANS(PLAN_NAME, CREATED_AT)
		VALUES('empty', NOW())	
		RETURNING ID
	`).Scan(&planId)
	if err != nil {
		return err
	}

	_, err = trnx.Exec(ctx, `
		INSERT INTO USER_PLANS(USER_ID, PLAN_ID)
		VALUES($1, $2)
	`, userId, planId)
	if err != nil {
		return err
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return err
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

func (d *DBs) GetAllUsersPlanIds(ctx context.Context, userId int) (*[]int, error) {

	var exerciseIds []int

	rows, err := d.PDB.Query(ctx, `
		SELECT PLAN_ID FROM USER_PLANS
		WHERE USER_ID = $1
	`, userId)

	if err != nil {
		return &exerciseIds, err
	}

	var id int

	for rows.Next() {

		err := rows.Scan(&id)
		if err != nil {
			return &exerciseIds, err
		}

		exerciseIds = append(exerciseIds, id)
	}

	rows.Close()

	return &exerciseIds, nil
}

func (d *DBs) GetEmptyPlanID(ctx context.Context, userId int) (int, error) {
	var planId int

	err := d.PDB.QueryRow(ctx, `
		SELECT PLANS.ID FROM USER_PLANS
		INNER JOIN PLANS
		ON USER_PLANS.PLAN_ID = PLANS.ID
		WHERE PLAN_NAME = 'empty' AND USER_ID = $1
	`, userId).Scan(&planId)

	if err != nil {
		return planId, err
	}

	return planId, nil
}

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

func (d *DBs) NoOfWorkoutsExistsInP() {

}

func (d *DBs) PlanExists(ctx context.Context, userId int, planName string) (bool, error) {
	var planId int
	err := d.PDB.QueryRow(ctx, `
		select plan_id from user_plans
		inner join plans
		on user_plans.plan_id = plans.id
		where user_id = $1 and plan_name = $2
	`, userId, planName).Scan(&planId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if the plan already exists : %w\n", err)
	}

	return true, nil
}

func (d *DBs) CreatePlan(ctx context.Context, userId int, p *models.Plan2) error {

	var planId int

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaction : %w\n", err)
	}

	defer trnx.Rollback(ctx)

	err = trnx.QueryRow(ctx, `	
		INSERT INTO PLANS(PLAN_NAME, CREATED_AT)
		VALUES($1, NOW())	
		RETURNING ID
	`, p.PlanName).Scan(&planId)
	if err != nil {
		return fmt.Errorf("error inserting plan name into plans : %w\n", err)
	}

	_, err = trnx.Exec(ctx, `
		INSERT INTO USER_PLANS(USER_ID, PLAN_ID)
		VALUES($1, $2)	
	`, userId, planId)
	if err != nil {
		return fmt.Errorf("error inserting userId and planId into user_plans :%w\n", err)
	}

	for _, v := range p.Exercises {
		eID, err := d.GetExerciseIdFromMain(ctx, v)
		if err != nil {
			return fmt.Errorf("error getting exercise id for %v from redis : %w\n", v, err)
		}
		_, err = trnx.Exec(ctx, `
			INSERT INTO PLAN_EXERCISES(PLAN_ID, EXERCISE_ID, REST_TIME_IN_SECONDS)
			VALUES($1, $2, 120)		
		`, planId, eID)
		if err != nil {
			return fmt.Errorf("error inserting exercise_id %v into plan_exercises : %w", eID, err)
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting the trnx for creating plan : %w", err)
	}

	return nil
}
