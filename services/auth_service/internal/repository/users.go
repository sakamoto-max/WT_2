package repository

import (
	"auth_service/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/sakamoto-max/wt_2-pkg/enum"
	myErrs "github.com/sakamoto-max/wt_2-pkg/my_errors"
)



func (r *repo) CreateUser(ctx context.Context, name string, email string, hashedPass string, role string) (string, time.Time, error) {
	
	trnx, err := r.pDB.Begin(ctx)
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
		"name" : name, 
		"email" : email, 
		"role" : role, 
		"hashedPass" : hashedPass,
	}).Scan(&userId, &createdAt)


	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_name_key":
				return "", createdAt, myErrs.AlreadyExitsErrMaker(string(enum.UserResource))
			case "users_email_key":
				return "", createdAt, myErrs.AlreadyExitsErrMaker(string(enum.UserResource))
			}
		}
		return "", createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("error commiting the transaction : %w", err))
	}

	dataForPlan := domain.EmptyPayload{
		UserId: userId,
	}

	planDataInBytes, err := json.Marshal(dataForPlan)
	if err != nil {
		return "", createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("failed to marshal the data : %w", err))
	}

	planPayload := string(planDataInBytes)
	
	dataForEmail := domain.EmailPayload{
		Email: email,
	}
	
	emailDataInBytes, err := json.Marshal(dataForEmail)
	if err != nil {
		return "", createdAt, myErrs.InternalServerErrMaker(fmt.Errorf("failed to marshal the data : %w", err))
	}

	emailPayload := string(emailDataInBytes)

	query = `
		INSERT INTO 
			outbox(target_service, created_by, task, payload)
		VALUES 
			(@planService, @createdBy, @emptyPlan, @planPayload::JSONB),
			(@emailService, @createdBy, @sendEmail, @emailPayload::JSONB)
	`

	_, err = trnx.Exec(ctx, query,pgx.NamedArgs{
		"planService" : enum.PlanService, 
		"emptyPlan" : enum.CreateEmptyPlanForUser, 
		"planPayload":planPayload,
		"emailService":enum.EmailService,
		"sendEmail" : enum.SendEmailforSigningUp,
		"emailPayload" : emailPayload,
		"createdBy" : "auth_service",
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




// DONE
func (r *repo) FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error) {

	query := `
		SELECT 
			ID, 
			NAME, 
			ROLE_ID 
		FROM 
			USERS
		WHERE 
			EMAIL = @email	
	`

	var userID string
	var roleID string
	var name string

	err := r.pDB.QueryRow(ctx, query, pgx.NamedArgs{"email" : email}).Scan(&userID, &name, &roleID)

	if err != nil {
		return userID, roleID, name, myErrs.InternalServerErrMaker(fmt.Errorf("error getting id, name, role_id : %w", err))
	}

	return userID, roleID, name, nil
}

func (r *repo) UserLogout(ctx context.Context, userId string, uuid string) error {
	// del refresh
	refreshKey := fmt.Sprintf("%v_refresh", uuid)
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	pipe := r.rDB.Pipeline()

	pipe.Del(ctx, refreshKey)
	pipe.Del(ctx, uuidKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myErrs.InternalServerErrMaker(fmt.Errorf("error deleting the refresh token after logout : %w\n", err))
	}

	return nil
}

// DONE
func (r *repo) FetchUserPass(ctx context.Context, email string) (string, error) {

	var hashedPass string

	query := `
		SELECT HASHED_PASS FROM USERS
		WHERE EMAIL = @email
	`

	err := r.pDB.QueryRow(ctx, query, pgx.NamedArgs{"email" : email}).Scan(&hashedPass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myErrs.ResourceNotFoundErrMaker(string(enum.EmailResource))
		}
		return "", myErrs.InternalServerErrMaker(fmt.Errorf("error getting hashed pass : %w", err))
	}

	return hashedPass, nil
}

// DONE
func (r *repo) FetchUserPassById(ctx context.Context, userId string) (string, error) {

	var hashedPass string

	query := `
		SELECT 
			HASHED_PASS 
		FROM 
			USERS
		WHERE
			id = @id	
	`

	err := r.pDB.QueryRow(ctx, query, pgx.NamedArgs{"id" : userId}).Scan(&hashedPass)
	if err != nil {
		return "", myErrs.InternalServerErrMaker(fmt.Errorf("error getting hashed pass : %w", err))
	}

	return hashedPass, nil
}

// DONE
func (r *repo) ChangePass(ctx context.Context, userId string, newPass string) error {
	
	query := `
		UPDATE 
			USERS
		SET 
			hashed_pass = @hashedPass
		WHERE 
			ID = @id
	
	`
	
	_, err := r.pDB.Exec(ctx, query, pgx.NamedArgs{"hashedPass" : newPass, "id" : userId})
	if err != nil {
		return myErrs.InternalServerErrMaker(fmt.Errorf("error updating the password : %w", err))
	}

	return nil
}

func (r *repo) GetUserId() {
	
}