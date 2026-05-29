package repository

import (
	"auth_service/internal/mappings"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	myErrs "github.com/sakamoto-max/wt_2_pkg/myerrs"

	// "github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
)

var (
	userResource  = "user"
	emailResource = "email"
)

type authDb struct {
	pg *pgxpool.Pool
}

func NewAuthDb(pg *pgxpool.Pool) *authDb {
	return &authDb{pg: pg}
}

func (a *authDb) CreateUser(ctx context.Context, payload mappings.SignUp) (string, time.Time, error) {

	trnx, err := a.pg.Begin(ctx)
	if err != nil {
		return "", time.Now(), myErrs.InternalServerErrMaker(fmt.Errorf("error creating transaction : %w\n", err))
	}

	defer trnx.Rollback(ctx)

	query := `
	INSERT INTO 
	users(name, email, role_id, hashed_pass)
	VALUES
	(@name, @email, (select id from roles where role = @role), @hashedPass)
	RETURNING 
	id, created_at
	`
	var userId string
	var createdAt time.Time

	err = trnx.QueryRow(ctx, query,
		pgx.NamedArgs{
			"name":       payload.Name,
			"email":      payload.Email,
			"role":       payload.Role,
			"hashedPass": payload.Password,
		}).Scan(&userId, &createdAt)

	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_name_key":
				return "", createdAt, myErrs.AlreadyExitsErrMaker(userResource)
			case "users_email_key":
				return "", createdAt, myErrs.AlreadyExitsErrMaker(userResource)
			}
		}
		return "", createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("error commiting the transaction : %w", err))
	}

	dataForPlan := map[string]string{
		enum.QueueFields_USER_ID.String(): userId,
	}

	planDataInBytes, err := json.Marshal(dataForPlan)
	if err != nil {
		return "", createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("failed to marshal the data : %w", err))
	}

	planPayload := string(planDataInBytes)

	dataForEmail := map[string]string{
		enum.QueueFields_EMAIL.String() : payload.Email,
	}

	emailDataInBytes, err := json.Marshal(dataForEmail)
	if err != nil {
		return "", createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("failed to marshal the data : %w", err))
	}

	emailPayload := string(emailDataInBytes)

	query = `
		INSERT INTO 
			outbox(target_service, created_by, task, status, payload)
		VALUES 
			(@planService, @createdBy, @emptyPlan, @status, @planPayload::JSONB),
			(@emailService, @createdBy, @sendEmail, @status, @emailPayload::JSONB)
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"planService":  enum.ServiceName_PLAN_SERVICE.String(),
		"emptyPlan":    enum.TaskName_CREATE_EMPTY_PLAN_FOR_USER.String(),
		"planPayload":  planPayload,
		"emailService": enum.ServiceName_EMAIL_SERVICE.String(),
		"sendEmail":    enum.TaskName_SEND_EMAIL_FOR_SIGNING_UP.String(),
		"emailPayload": emailPayload,
		"createdBy":    enum.ServiceName_AUTH_SERVICE.String(),
		"status":       enum.TaskStatus_TASK_NOT_COMPLETED.String(),
	})

	if err != nil {
		return userId, createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("error inserting data into outbox : %w", err))
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return userId, createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("error commiting the transaction : %w", err))
	}

	return userId, createdAt, nil
}
