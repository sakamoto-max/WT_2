package repository

import (
	"context"
	"errors"
	"fmt"
	"plan_service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type planQueryRepo struct {
	pg *pgxpool.Pool
}

func (d *planQueryRepo) GetAllPlanNamesWithIds(ctx context.Context, userId string) (*[]domain.Plan, error) {

	var allPlans []domain.Plan

	query := `
		SELECT 
			ID, 
			NAME 
		FROM 
			PLANS
		WHERE 
			USER_ID = @userId
		`
	rows, err := d.pg.Query(ctx, query, pgx.NamedArgs{"userId": userId})

	if err != nil {
		return nil, myerrors.InternalServerErrMaker(fmt.Errorf("error getting plan ids for the user %v : %w", userId, err))
	}

	var planId string
	var planName string

	for rows.Next() {

		err := rows.Scan(&planId, &planName)
		if err != nil {
			return &allPlans, fmt.Errorf("error scaning rows : %w", err)
		}

		a := domain.Plan{PlanName: planName, PlanId: planId}

		allPlans = append(allPlans, a)
	}

	rows.Close()
	return &allPlans, nil
}

func (d *planQueryRepo) GetEmptyPlanId(ctx context.Context, userId string) (string, error) {
	var planId string
	query := `
		SELECT
			id 
		FROM
			plans
		WHERE 
			user_id = @userId AND NAME = @name
	`

	err := d.pg.QueryRow(ctx, query, pgx.NamedArgs{"userId": userId, "name": "empty"}).Scan(&planId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return planId, myerrors.ResourceNotFoundErrMaker("plan")
		}
		return planId, myerrors.InternalServerErrMaker(fmt.Errorf("error checking if the plan already exists : %w\n", err))
	}

	return planId, nil

}

func (d *planQueryRepo) GetPlanId(ctx context.Context, payload domain.GetPlan) (string, error) {
	var planId string

	query := `
		SELECT
			id 
		FROM
			plans
		WHERE 
			user_id = @userId AND NAME = @name
	`

	err := d.pg.QueryRow(ctx, query, pgx.NamedArgs{"userId": payload.UserId, "name": payload.PlanName}).Scan(&planId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return planId, myerrors.ResourceNotFoundErrMaker("plan")
		}
		return planId, myerrors.InternalServerErrMaker(fmt.Errorf("error checking if the plan already exists : %w\n", err))
	}

	return planId, nil
}

func (d *planQueryRepo) GetAllExercisesByPlanID(ctx context.Context, planId string) (*[]string, error) {
	var exerciseIDs []string

	query := `
		SELECT 	
			EXERCISE_ID 
		FROM 
			PLAN_EXERCISES
		WHERE 
			PLAN_ID = @planId
	`

	rows, err := d.pg.Query(ctx, query, pgx.NamedArgs{"planId": planId})
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
func (d *planQueryRepo) GetPlan(ctx context.Context, payload domain.GetPlan) (string, *[]string, error) {
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

	rows, err := d.pg.Query(ctx, query, pgx.NamedArgs{
		"userId":   payload.UserId,
		"planName": payload.PlanName,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
