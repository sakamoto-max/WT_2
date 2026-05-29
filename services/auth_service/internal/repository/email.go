package repository

import (
	"context"
	"fmt"

	// myErrs "github.com/sakamoto-max/wt_2-pkg/my_errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myErrs "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type emailDb struct {
	pg *pgxpool.Pool
}

type EmailIface interface {
	GetEmail(ctx context.Context, userId string) (string, error)
	ChangeEmail(ctx context.Context, userId string, newEmail string) error
}

func NewEmailDb(pg *pgxpool.Pool) EmailIface {
	return &emailDb{pg: pg}
}

func (e *emailDb) GetEmail(ctx context.Context, userId string) (string, error) {

	var email string

	query := `
		SELECT 
			email 
		FROM 
			users
		WHERE 
			id = @id
	`

	err := e.pg.QueryRow(ctx, query, pgx.NamedArgs{"id": userId}).Scan(&email)
	if err != nil {
		return email, myErrs.InternalServerErrMaker(fmt.Errorf("error getting email of user with id : %v : %w", userId, err))
	}

	return email, nil
}
func (e *emailDb) ChangeEmail(ctx context.Context, userId string, newEmail string) error {

	query := `
		UPDATE users
		SET email = @email
		WHERE id = @id	
	`
	_, err := e.pg.Exec(ctx, query, pgx.NamedArgs{"email": newEmail, "id": userId})
	if err != nil {
		return myErrs.InternalServerErrMaker(fmt.Errorf("error changing the email in the db : %w", err))
	}

	return nil
}
