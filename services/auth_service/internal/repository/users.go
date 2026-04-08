package repository

import (
	// customerrors "auth_service/internal/custom_errors"
	"auth_service/internal/models"
	"auth_service/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"
	"wt/pkg/enum"
	myErrs "wt/pkg/my_errors"

	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
)


// DONE
func (r *Repo) CreateUser(ctx context.Context, name string, email string, hashedPass string, role string) (string, time.Time, error) {
	
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

	dataForPlan := models.EmptyPayload{
		UserId: userId,
	}

	dataForEmail := models.EmailPayload{
		Email: email,
	}

	payloadForPlan, _ := utils.MakeJSONV2(dataForPlan)
	payloadForEmail, _ := utils.MakeJSONV2(dataForEmail)

	query = `
		INSERT INTO 
			outbox(target_service, task, payload)
		VALUES 
			(@planService, @emptyPlan, @planPayload::JSONB),
			(@emailService, @sendEmail, @emailPayload::JSONB)
	`

	_, err = trnx.Exec(ctx, query,pgx.NamedArgs{
		"planService" : enum.PlanService, 
		"emptyPlan" : enum.CreateEmptyPlanForUser, 
		"planPayload":payloadForPlan,
		"emailService":enum.EmailService,
		"sendEmail" : enum.SendEmailforSigningUp,
		"emailPayload" : payloadForEmail,
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
func (r *Repo) FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error) {

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

func (r *Repo) UserLogout(ctx context.Context, userId string, uuid string) error {
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
func (r *Repo) FetchUserPass(ctx context.Context, email string) (string, error) {

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
func (r *Repo) FetchUserPassById(ctx context.Context, userId string) (string, error) {

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
func (r *Repo) ChangePass(ctx context.Context, userId string, newPass string) error {
	
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

func (r *Repo) GetUserId() {
	
}