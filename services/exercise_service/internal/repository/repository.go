package repository

import (
	"context"
	"errors"
	"exercise_service/internal/domain"
	"fmt"
	"time"
	enum "wt/pkg/enum"
	myerrs "wt/pkg/my_errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) GetPostgresRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.pDB.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *Repo) GetRedisRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.rDB.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *Repo) GetExerciseByName(ctx context.Context, userId string, exerciseName string) (*domain.Exercise, error) {

	var exercise domain.Exercise
	query := `
			SELECT 
			EXERCISES.ID, 
			EXERCISES.NAME, 
			BODY_PARTS.NAME, 
			EQUIPMENT.NAME ,
			created_at,
			updated_at
		FROM 
			EXERCISES
		INNER JOIN 
			BODY_PARTS 
		ON 
			EXERCISES.BODY_PART_ID = BODY_PARTS.ID
		INNER JOIN 
			EQUIPMENT 
		ON 
			EXERCISES.EQUIPMENT_ID = EQUIPMENT.ID
		FULL JOIN (
			SELECT * FROM USER_NULLIFIED
			WHERE USER_ID = @userId
			) AS TABLE_B
		ON 
			EXERCISES.ID = TABLE_B.EXERCISE_ID
		WHERE 
			Exercises.NAME = @exerciseName and (CREATED_BY IS NULL OR CREATED_BY = @userId) AND TABLE_B.EXERCISE_ID IS NULL;	
	`
	row := r.pDB.QueryRow(ctx, query, pgx.NamedArgs{"exerciseName": exerciseName, "userId": userId})
	err := row.Scan(&exercise.Id, &exercise.Name, &exercise.BodyPart, &exercise.Equipment, &exercise.CreatedAt, &exercise.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrs.ResourceNotFoundErrMaker(string(enum.ExerciseResource))
		}

		return nil, myerrs.InternalServerErrMaker(fmt.Errorf("error in getting the exercise by name : %v : %w", exerciseName, err))
	}

	return &exercise, nil
}

var (
	id        string = "id"
	bodyPart  string = "body_part"
	equipment string = "equipment"
	createdAt string = "created_at"
	updatedAt string = "updated_at"
)

func (r *Repo) GetExerciseByNameR(ctx context.Context, userId string, exerciseName string) (*domain.Exercise, error) {
	
	key := fmt.Sprintf("user_id:%v:exercise_name:%v", userId, exerciseName)

	res, err := r.rDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise by name from cache : %w", err)
	}

	if len(res) == 0 {
		return nil, nil
	}

	layout := "2006-01-02T15:04:05.9999999Z07:00"
	createdAt, err := time.Parse(layout, res[createdAt])
	updatedAt, err := time.Parse(layout, res[updatedAt])

	id := res[id]
	name := exerciseName
	bodyPart := res[bodyPart]
	equipment := res[equipment]

	data := domain.Exercise{
		Id:        id,
		Name:      name,
		BodyPart:  bodyPart,
		Equipment: equipment,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return &data, nil
}

func (r *Repo) SetExerciseByNameR(ctx context.Context, userId string, exerData *domain.Exercise) error {
	mainKey := fmt.Sprintf("user_id:%v:exercise_name:%v", userId, exerData.Name)
	idKey := "id"
	bodyPartKey := "body_part"
	equipmentKey := "equipment"
	createdAtKey := "created_at"
	updatedAtKey := "updated_at"

	err := r.rDB.HSet(ctx, mainKey,
		idKey, exerData.Id,
		bodyPartKey, exerData.BodyPart,
		equipmentKey, exerData.Equipment,
		createdAtKey, exerData.CreatedAt,
		updatedAtKey, exerData.UpdatedAt,
	).Err()

	if err != nil {
		return fmt.Errorf("error setting exercise %v in redis : %w", exerData.Name, err)
	}

	return nil
}

func (r *Repo) CreateExercise(ctx context.Context, userId string, exerciseName string, bodyPartName string, equipmentName string) (string, error) {

	var bodyPartId uuid.UUID
	var equipmentId uuid.UUID
	var Id uuid.UUID

	trnx, err := r.pDB.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating a transaction : %w\n", err)
	}

	defer trnx.Rollback(ctx)

	err = trnx.QueryRow(ctx, `
		SELECT ID FROM BODY_PARTS
		WHERE NAME = @name
	`, pgx.NamedArgs{"name": bodyPartName}).Scan(&bodyPartId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myerrs.ResourceNotFoundErrMaker(string(enum.BodyPartResource))
		}

		err := fmt.Errorf("error getting id of body_part : %w\n", err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	err = trnx.QueryRow(ctx, `
		SELECT ID FROM EQUIPMENT
		WHERE NAME = @name
	`, pgx.NamedArgs{"name": equipmentName}).Scan(&equipmentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myerrs.ResourceNotFoundErrMaker(string(enum.EquipmentResource))
		}

		err := fmt.Errorf("error getting id of equipment : %w\n", err)
		return "", myerrs.InternalServerErrMaker(err)
	}
	err = trnx.QueryRow(ctx, `
		INSERT INTO exercises(name, created_by, body_part_id, equipment_id)
		VALUES(	@name, @createdBy, @bodyPartId, @equipmentId)
		RETURNING ID
	`, pgx.NamedArgs{"name": exerciseName, "createdBy": userId, "bodyPartId": bodyPartId, "equipmentId": equipmentId}).Scan(&Id)
	if err != nil {
		err := fmt.Errorf("error inserting exercise %v : %w\n", exerciseName, err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		err := fmt.Errorf("error commiting the transaction : %w\n", err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	return Id.String(), nil
}
func (r *Repo) GetAllExercises(ctx context.Context, userId string) (*[]domain.Exercise, error) {
	// HSET ALLEXERCISES EXERCISE_NAME

	query := `
		SELECT 
			EXERCISES.ID, 
			EXERCISES.NAME, 
			BODY_PARTS.NAME, 
			EQUIPMENT.NAME ,
			created_at,
			updated_at
		FROM 
			EXERCISES
		INNER JOIN 
			BODY_PARTS 
		ON 
			EXERCISES.BODY_PART_ID = BODY_PARTS.ID
		INNER JOIN 
			EQUIPMENT 
		ON 
			EXERCISES.EQUIPMENT_ID = EQUIPMENT.ID
		FULL JOIN (
			SELECT * FROM USER_NULLIFIED
			WHERE USER_ID = @userId
			) AS TABLE_B
		ON 
			EXERCISES.ID = TABLE_B.EXERCISE_ID
		WHERE 
			(CREATED_BY IS NULL OR CREATED_BY = @userId) AND TABLE_B.EXERCISE_ID IS NULL;
	`

	var allExercises []domain.Exercise

	rows, err := r.pDB.Query(ctx, query, pgx.NamedArgs{"userId": userId})
	if err != nil {

		err := fmt.Errorf("error getting all exercises from DB : %w\n", err)

		return nil, myerrs.InternalServerErrMaker(err)
	}

	var id string
	var exerciseName string
	var bodyPart string
	var equipmentName string
	var createdAt time.Time
	var updatedAt time.Time

	for rows.Next() {
		err := rows.Scan(&id, &exerciseName, &bodyPart, &equipmentName, &createdAt, &updatedAt)
		if err != nil {
			err := fmt.Errorf("error scanning the rows : %w\n", err)
			return nil, myerrs.InternalServerErrMaker(err)
		}

		exercise := domain.Exercise{
			Id:        id,
			Name:      exerciseName,
			BodyPart:  bodyPart,
			Equipment: equipmentName,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		allExercises = append(allExercises, exercise)
	}

	return &allExercises, nil
}
func (r *Repo) GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error) {

	var exerciseName string

	err := r.pDB.QueryRow(ctx, `
		SELECT name from exercises
		where id = @id	
	`, pgx.NamedArgs{"id": exerciseId}).Scan(&exerciseName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myerrs.ResourceNotFoundErrMaker(string(enum.ExerciseResource))
		}

		err := fmt.Errorf("error getting execise name for id %v : %w", exerciseId, err)
		return "", myerrs.InternalServerErrMaker(err)
	}

	return exerciseName, nil
}
func (r *Repo) DeleteExecise(ctx context.Context, userId string, exerciseName string) error {
	// if the exercise's created by is NULL -> move it to user nullified
	// else delete the exercise

	trnx, err := r.pDB.Begin(ctx)
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error creating a transaction"))
	}

	defer trnx.Rollback(ctx)

	query := `SELECT 
				id, 
				created_by 
			FROM
				exercises
			WHERE 
				name = @name AND (created_by IS NULL OR created_by = @userId) `

	var exerciseId uuid.UUID
	var createdBy *int
	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{"name": exerciseName, "userId": userId}).Scan(&exerciseId, &createdBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return myerrs.ResourceNotFoundErrMaker(string(enum.ExerciseResource))
		}
		return myerrs.InternalServerErrMaker(fmt.Errorf("error getting id and createdBy for exericse %v : %w", exerciseName, err))
	}

	if createdBy == nil {
		query := `
			INSERT INTO user_nullified(user_id, exercise_id)
			VALUES(@userId, @exerciseId)	
		`
		_, err := trnx.Exec(ctx, query, pgx.NamedArgs{"userId": userId, "exerciseId": exerciseId})
		if err != nil {
			return myerrs.InternalServerErrMaker(fmt.Errorf("error inserting data into user_nullified : %w", err))
		}

		err = trnx.Commit(ctx)
		if err != nil {
			return myerrs.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
		}

		return nil
	}

	query = `
		DELETE FROM exercises
		WHERE id = @id
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{"id": exerciseId})
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error deleting rows from table : %w", err))
	}
	err = trnx.Commit(ctx)
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("error committing : %w", err))
	}

	return nil
}
func (r *Repo) ExerciseExistsReturnId(ctx context.Context, userId string, exerciseName string) (string, error) {

	query := `
		SELECT 
			ID 
		FROM 
			EXERCISES
		LEFT JOIN 
			(SELECT EXERCISE_ID FROM USER_NULLIFIED WHERE USER_ID = @userId) AS TABLE_2
		ON 
			EXERCISES.ID = TABLE_2.EXERCISE_ID
		WHERE 
			NAME = @exerciseName
		AND 
			(CREATED_BY = @userId OR CREATED_BY IS NULL) 
		AND 
			TABLE_2.EXERCISE_ID IS NULL;	
	`
	var id string
	err := r.pDB.QueryRow(ctx, query, pgx.NamedArgs{"exerciseName": exerciseName, "userId": userId}).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return id, myerrs.ResourceNotFoundErrMaker(exerciseName)
		}
		return id, myerrs.InternalServerErrMaker(fmt.Errorf("error checking if the exericse exists : %v", err))
	}

	return id, nil
}
