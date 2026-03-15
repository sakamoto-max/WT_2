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
