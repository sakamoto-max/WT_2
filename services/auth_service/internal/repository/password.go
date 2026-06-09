package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myErrs "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type passDb struct {
	pg *pgxpool.Pool
}

type PassDbIface interface {
	FetchUserPass(ctx context.Context, email string) (string, error) 
	FetchUserPassById(ctx context.Context, userId string) (string, error)
	ChangePass(ctx context.Context, userId string, newPass string) error
}

func NewPassDb(pg *pgxpool.Pool) PassDbIface {
	return &passDb{pg : pg}
}

func (p *passDb) FetchUserPass(ctx context.Context, email string) (string, error) {

	var hashedPass string

	query := `
		SELECT HASHED_PASS FROM USERS
		WHERE EMAIL = @email
	`
	err := p.pg.QueryRow(ctx, query, pgx.NamedArgs{"email": email}).Scan(&hashedPass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", myErrs.ResourceNotFoundErrMaker(string(userResource))
		}
		return "", myErrs.InternalServerErrMaker(fmt.Errorf("error getting hashed pass : %w", err))
	}

	return hashedPass, nil
}
func (p *passDb) FetchUserPassById(ctx context.Context, userId string) (string, error) {

	var hashedPass string

	query := `
		SELECT 
			HASHED_PASS 
		FROM 
			USERS
		WHERE
			id = @id	
	`

	err := p.pg.QueryRow(ctx, query, pgx.NamedArgs{"id": userId}).Scan(&hashedPass)
	if err != nil {
		return "", myErrs.InternalServerErrMaker(fmt.Errorf("error getting hashed pass : %w", err))
	}

	return hashedPass, nil
}
func (p *passDb) ChangePass(ctx context.Context, userId string, newPass string) error {

	query := `
		UPDATE 
			USERS
		SET 
			hashed_pass = @hashedPass
		WHERE 
			ID = @id
	
	`

	_, err := p.pg.Exec(ctx, query, pgx.NamedArgs{"hashedPass": newPass, "id": userId})
	if err != nil {
		return myErrs.InternalServerErrMaker(fmt.Errorf("error updating the password : %w", err))
	}

	return nil
}
