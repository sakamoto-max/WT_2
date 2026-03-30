package repository

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/models"
	"wt/pkg/enum"
	myerrors "wt/pkg/my_errors"

	// "plan_service/internal/models"

	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrPlanAlreadyExists = errors.New("plan already exits")
	// ErrNoExerciseExits = errors.New("no exercise ext")
)

func (d *DBs) CreatePlan(ctx context.Context, userId int, planName string, exerciseIds []string) error {

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating transaction : %w\n", err))
	}

	defer trnx.Rollback(ctx)

	query := `
		INSERT INTO PLANS(USER_ID, NAME)
		VALUES (@userId, @name)	
		RETURNING ID
	`
	var planId string

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"userId": userId, "name": planName}).Scan(&planId)
	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "one_user_one_plan_name" {
			return myerrors.BadReqErrMaker(ErrPlanAlreadyExists)
		}
		return myerrors.InternalServerErrMaker(fmt.Errorf("error inserting plan name into plans : %w\n", err))
	}

	query = `
		INSERT INTO PLAN_EXERCISES(PLAN_ID, EXERCISE_ID)
		VALUES(@planId, @exerciseId)	
	`

	for _, exerciseId := range exerciseIds {
		fmt.Println(exerciseId)
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId, "exerciseId": exerciseId})
		if err != nil {
			var pgErr *pgConn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "one_plan_id_one_exercise" {
				return myerrors.BadReqErrMaker(fmt.Errorf("exercise %v is selected twice", exerciseId))
			}
			return myerrors.InternalServerErrMaker(fmt.Errorf("error inserting exercise_id %v into plan_exercises : %w", exerciseId, err))
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error commiting the trnx for creating plan : %w", err))
	}

	return nil
}

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
func (d *DBs) ReturnsPlanId(ctx context.Context, userId int, planName string) (string, error) {
	var planId string

	query := `
		SELECT
			id 
		FROM
			plans
		WHERE 
			user_id = @userId AND NAME = @name
	`

	err := d.PDB.QueryRow(ctx, query, pgx.NamedArgs{"userId": userId, "name": planName}).Scan(&planId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return planId, myerrors.ResourceNotFoundErrMaker(string(enum.PlanResource))
		}
		return planId, myerrors.InternalServerErrMaker(fmt.Errorf("error checking if the plan already exists : %w\n", err))
	}

	return planId, nil
}

func (d *DBs) GetAllExercisesByPlanID(ctx context.Context, planId string) (*[]string, error) {
	var exerciseIDs []string

	query := `
		SELECT 	
			EXERCISE_ID 
		FROM 
			PLAN_EXERCISES
		WHERE 
			PLAN_ID = @planId
	`

	rows, err := d.PDB.Query(ctx, query, pgx.NamedArgs{"planId": planId})
	if err != nil {
		return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting exercises for the plan %v : %w", planId, err))
	}

	defer rows.Close()

	var id string

	for rows.Next() {

		err := rows.Scan(&id)
		if err != nil {
			return &exerciseIDs, myerrors.InternalServerErrMaker(fmt.Errorf("error in scaning the rows : %w", err))
		}

		exerciseIDs = append(exerciseIDs, id)
	}

	return &exerciseIDs, nil
}

// DONE
func (d *DBs) AddExercisesToPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating a transaction : %w", err))
	}

	defer trnx.Rollback(ctx)

	query := `
			INSERT INTO 
				plan_exercises(plan_id, exercise_id)
			VALUES 
				(@planId, @exerciseId)
		`

	for _, id := range *exerciseIDs {
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"planId" : planId, "exerciseId" : id})
		if err != nil {
			return myerrors.InternalServerErrMaker(fmt.Errorf("error inserting the exericse with id %v : %w", id, err)) 
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil
}

// DONE
func (d *DBs) DeleteExerciseFromPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {

	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating a transaction : %w", err))
	}

	defer trnx.Rollback(ctx)

	query := `
		DELETE FROM 
			plan_exercises
		WHERE 
			PLAN_ID = @planId AND EXERCISE_ID = @exerciseId
	`

	for _, id := range *exerciseIDs {
		_, err := trnx.Exec(ctx, query,pgx.NamedArgs{"planId" : planId, "exerciseId" : id})
		if err != nil {
			return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting exercise with id %v : %w", id, err))
		}
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil
}
func (d *DBs) GetAllUserPlans(ctx context.Context, userId int) (*[]models.Plan3, error) {

	var allPlans []models.Plan3

	query := `
		SELECT 
			ID, 
			NAME 
		FROM 
			PLANS
		WHERE 
			USER_ID = @userId
		`
	rows, err := d.PDB.Query(ctx, query, pgx.NamedArgs{"userId": userId})
	defer rows.Close()

	if err != nil {
		return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting plan ids for the user %v : %w", userId, err))
	}

	var id string
	var planName string

	for rows.Next() {

		err := rows.Scan(&id, &planName)
		if err != nil {
			return &allPlans, fmt.Errorf("error scaning rows : %w", err)
		}

		a := models.Plan3{PlanName: planName, Id: id}

		allPlans = append(allPlans, a)
	}

	return &allPlans, nil
}

// DONE
func (d *DBs) DeletePlan(ctx context.Context, userId int, planId string) error {
	trnx, err := d.PDB.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating a transaction : %w", err))
	}

	query := `
		DELETE FROM 
			plan_exercises
		WHERE 
			PLAN_ID = @planId
	`

	_, err = trnx.Exec(ctx,query, pgx.NamedArgs{"planId" : planId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting exercises : %w", err))
	}

	query = `
		DELETE FROM 
			PLANS
		WHERE 
			ID = @id
	`

	_, err = trnx.Exec(ctx, query,pgx.NamedArgs{"id" : planId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting plan : %w", err))
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil

}

func (d *DBs) GetEmptyPlanID(ctx context.Context, userId int) (int, error) {
	var emptyPlanId int
	err := d.PDB.QueryRow(ctx, `
		SELECT id FROM plans
		WHERE name = 'empty' AND user_id = $1;
	`, userId).Scan(&emptyPlanId)
	if err != nil {
		return emptyPlanId, fmt.Errorf("error getting empty plan Id : %w", err)
	}

	return emptyPlanId, nil

}

// func (d *DBs) ReturnPlan(ctx context.Context, userId int, planName string) {

// 	query := `
// 		SELECT 
// 	`

// }

// func (d *DBs) PlanExistsReturnPlan(ctx context.Context, userID int, planName string) (bool, int, *[]int, error) {

// 	var allExerciseIDs *[]int
// 	exists, planId, err := d.(ctx, userID, planName)
// 	if err != nil {
// 		return exists, planId, allExerciseIDs, fmt.Errorf("error getting exercise ids : %w", err)
// 	}
// 	if !exists {
// 		return exists, planId, allExerciseIDs, nil
// 	}

// 	allExerciseIDs, err = d.GetAllExercisesByPlanID(ctx, planId)
// 	if err != nil {
// 		return exists, planId, allExerciseIDs, fmt.Errorf("error getting exercise ids : %w", err)
// 	}

// 	return exists, planId, allExerciseIDs, nil
// }

func (d *DBs) CreateEmptyPlan(ctx context.Context, userId int) error {

	query := `
		INSERT INTO plans(user_id, name, created_at)
		VALUES(@userId, @name, NOW())
	`
	_, err := d.PDB.Exec(ctx, query, pgx.NamedArgs{"userId": userId, "name": models.EmptyPlan})
	if err != nil {
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
