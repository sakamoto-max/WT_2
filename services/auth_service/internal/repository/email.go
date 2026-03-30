package repository

import (
	"context"
	"fmt"

	myErrs "wt/pkg/my_errors"

	"github.com/jackc/pgx/v5"
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

// DONE
func (r *Repo) GetEmail(ctx context.Context, userId string) (string, error) {
	
	var email string

	query :=  `
		SELECT 
			email 
		FROM 
			users
		WHERE 
			id = @id
	`

	err := r.pDB.QueryRow(ctx, query,pgx.NamedArgs{"id" : userId}).Scan(&email)
	if err != nil{
		return email, myErrs.InternalServerErrMaker(fmt.Errorf("error getting email of user with id : %v : %w", userId, err))
	}

	return email, nil
}


func (r *Repo) ChangeEmail(ctx context.Context, userId string, newEmail string) error {

	query := `
		UPDATE users
		SET email = @email
		WHERE id = @id	
	`
	_, err := r.pDB.Exec(ctx, query, pgx.NamedArgs{"email" : newEmail, "id" : userId})
	if err != nil{
		return myErrs.InternalServerErrMaker(fmt.Errorf("error changing the email in the db : %w",err))
	}

	return nil
}