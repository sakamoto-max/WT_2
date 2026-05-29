package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myErrs "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type UserDbIface interface {
	FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error)
}
type userDb struct {
	pg *pgxpool.Pool
}

func NewUserDb(pg *pgxpool.Pool) UserDbIface {
	return &userDb{pg: pg}
}

func (u *userDb) FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error) {

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

	err := u.pg.QueryRow(ctx, query, pgx.NamedArgs{"email": email}).Scan(&userID, &name, &roleID)

	if err != nil {
		return userID, roleID, name, myErrs.InternalServerErrMaker(fmt.Errorf("error getting id, name, role_id : %w", err))
	}

	return userID, roleID, name, nil
}
