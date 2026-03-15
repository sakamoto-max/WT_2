package repository

import (
	// customerrors "auth_service/internal/custom_errors"
	"context"
	"errors"
	"fmt"
	"time"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	myErrs "wt/pkg/my_errors"
)

func (r *Repo) CreateUser(ctx context.Context, name string, email string, hashedPass string) (int, time.Time, error) {
	var userId int
	var createdAt time.Time

	err := r.pDB.QueryRow(ctx, `
	INSERT INTO users(name, email, role_id, hashed_pass, created_at)
	VALUES($1, $2, (select id from roles where role = 'user'), $3, NOW())
	RETURNING id, created_at
	`, name, email, hashedPass).Scan(&userId, &createdAt)

	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_name_key":
				return userId, createdAt, myErrs.ErrNameAlreadyExits
			case "users_email_key":
				return userId, createdAt, myErrs.ErrEmailAlreadyExits
			}
		}
		return userId, createdAt, err
	}

	return userId, createdAt, nil
}

func (r *Repo) FetchNameUserIdRoleId(ctx context.Context, email string) (int, int, string, error) {

	var userID int
	var roleID int
	var name string

	err := r.pDB.QueryRow(ctx, `
		SELECT ID, NAME, ROLE_ID FROM USERS
		WHERE EMAIL = $1
	`, email).Scan(&userID, &name, &roleID)

	if err != nil {
		return userID, roleID, name, fmt.Errorf("error getting id, name, role_id : %w", err)
	}

	return userID, roleID, name, nil
}

func (r *Repo) UserLogout(ctx context.Context, userId int, uuid string) error {
	// del refresh
	refreshKey := fmt.Sprintf("%v_refresh", uuid)
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	pipe := r.rDB.Pipeline()

	pipe.Del(ctx, refreshKey)
	pipe.Del(ctx, uuidKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting the refresh token after logout : %w\n", err)
	}

	return nil
}

func (r *Repo) FetchUserPass(ctx context.Context, email string) (string, error) {

	var hashedPass string

	err := r.pDB.QueryRow(ctx, `
		SELECT HASHED_PASS FROM USERS
		WHERE EMAIL = $1	
	`, email).Scan(&hashedPass)
	if err != nil {
		return "", fmt.Errorf("error getting hashed pass : %w", err)
	}

	return hashedPass, nil
}
