package repository

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

var (
	ErrPlanAlreadyExists = errors.New("plan already exits")
)

// NEED
func (d *dBs) CreatePlan(ctx context.Context, userId string, planName string, exerciseIds []string) error {

	trnx, err := d.pDB.Begin(ctx)
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
func (d *dBs) GetPlans(ctx context.Context, userId string) (*[]models.Plan3, error) {

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
	rows, err := d.pDB.Query(ctx, query, pgx.NamedArgs{"userId": userId})

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

	rows.Close()
	return &allPlans, nil
}
func (d *dBs) GetAllExercisesByPlanID(ctx context.Context, planId string) (*[]string, error) {
	var exerciseIDs []string

	query := `
		SELECT 	
			EXERCISE_ID 
		FROM 
			PLAN_EXERCISES
		WHERE 
			PLAN_ID = @planId
	`

	rows, err := d.pDB.Query(ctx, query, pgx.NamedArgs{"planId": planId})
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
func (d *dBs) ReturnsPlanId(ctx context.Context, userId string, planName string) (string, error) {
	var planId string

	query := `
		SELECT
			id 
		FROM
			plans
		WHERE 
			user_id = @userId AND NAME = @name
	`

	err := d.pDB.QueryRow(ctx, query, pgx.NamedArgs{"userId": userId, "name": planName}).Scan(&planId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return planId, myerrors.ResourceNotFoundErrMaker("plan")
		}
		return planId, myerrors.InternalServerErrMaker(fmt.Errorf("error checking if the plan already exists : %w\n", err))
	}

	return planId, nil
}
func (d *dBs) AddExercisesToPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {

	trnx, err := d.pDB.Begin(ctx)
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
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId, "exerciseId": id})
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
func (d *dBs) DeleteExerciseFromPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {

	trnx, err := d.pDB.Begin(ctx)
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
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId, "exerciseId": id})
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
func (d *dBs) DeletePlan(ctx context.Context, userId string, planId string) error {
	trnx, err := d.pDB.Begin(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error creating a transaction : %w", err))
	}

	query := `
		DELETE FROM 
			plan_exercises
		WHERE 
			PLAN_ID = @planId
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"planId": planId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting exercises : %w", err))
	}

	query = `
		DELETE FROM 
			PLANS
		WHERE 
			ID = @id
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"id": planId})
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting plan : %w", err))
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil

}
func (d *dBs) CreateEmptyPlan(ctx context.Context, userId string) error {

	query := `
		INSERT INTO plans(user_id, name)
		VALUES(@userId, @name)
	`
	_, err := d.pDB.Exec(ctx, query, pgx.NamedArgs{"userId": userId, "name": "empty"})
	if err != nil {
		return fmt.Errorf("Error creating empty plan : %v", err)
	}

	return nil
}
func (r *dBs) GetPostgresRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.pDB.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *dBs) GetRedisRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.rDB.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *dBs) GetPlan(ctx context.Context, userId string, planName string) (string, *[]string, error) {
	query := `
		SELECT 
			PLANS.ID, 
			EXERCISE_ID 
		FROM 
			PLANS
		INNER JOIN 
			PLAN_EXERCISES
		ON 
			PLANS.ID = PLAN_EXERCISES.PLAN_ID
		WHERE 
			user_id= @userId
		AND 
			name=@planName
	`

	rows, err := r.pDB.Query(ctx, query, pgx.NamedArgs{
		"userId" : userId,
		"planName": planName,
	})


	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return "", nil, myerrors.ResourceNotFoundErrMaker("plan")
		}

		return "", nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting plan details : %w", err))
	}

	var planId string
	var exerciseId string
	
	var exerciseIds []string
	// var planIdWithAllExericseIds []models.PlanIDExerciseId

	for rows.Next() {
		err := rows.Scan(&planId, &exerciseId)
		if err != nil {
			return "", nil, myerrors.InternalServerErrMaker(fmt.Errorf("error scanning the rows : %w", err))
		}
	
		exerciseIds = append(exerciseIds, exerciseId)
	}

	rows.Close()

	return planId, &exerciseIds, nil
}