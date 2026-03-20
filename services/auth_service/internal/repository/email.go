package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	myErrs "wt/pkg/my_errors"
)

func (r *Repo) EmailExists(ctx context.Context, email string) (bool, error) {

	var id int

	err := r.pDB.QueryRow(ctx, `
		SELECT ID FROM USERS
		WHERE EMAIL = $1	
	`, email).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, myErrs.ErrEmailNotFound
		}

		return false, fmt.Errorf("error checking if email %v exists : %w\n", email, err)
	}

	return true, nil
}

func (r *Repo) GetEmail(ctx context.Context, userId int) (string, error) {
	
	var email string
	
	err := r.pDB.QueryRow(ctx, `
		SELECT email FROM users
		WHERE id = $1	
	`, userId).Scan(&email)
	if err != nil{
		return email, fmt.Errorf("error getting email of user with id : %v : %w", userId, err)
	}

	return email, nil
}

func (r *Repo) ChangeEmail(ctx context.Context, userId int, newEmail string) error {
	_, err := r.pDB.Exec(ctx, `
		UPDATE users
		SET email = $1
		WHERE id = $2	
	`, newEmail, userId)
	if err != nil{
		return fmt.Errorf("error changing the email in the db : %w",err)
	}

	return nil
}