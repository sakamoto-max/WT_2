package repository

import (
	"auth_service/internal/responses"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repo struct {
	pDB *pgxpool.Pool
	rDB *redis.Client
}

func NewRepo(pool *pgxpool.Pool, client *redis.Client) *Repo {
	return &Repo{pDB: pool, rDB: client}
}

// postgres
func (r *Repo) SignUp(ctx context.Context, name string, email string, hashedPass string) (*responses.SignUpResp, error) {

	var userId int
	var resp responses.SignUpResp
	var roleId int

	trnx, err := r.pDB.Begin(ctx)
	if err != nil {
		return &resp, fmt.Errorf("error creating the transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	err = trnx.QueryRow(ctx, `
		INSERT INTO USERS(NAME, EMAIL, HASHED_PASS, CREATED_AT)
		VALUES($1, $2, $3, NOW())	
		RETURNING ID, NAME, EMAIL, CREATED_AT
	`, name, email, hashedPass).Scan(&userId, &resp.Name, &resp.Email, &resp.CreatedAt)
	if err != nil {
		return &resp, fmt.Errorf("error inserting data into users : %w", err)
	}

	err = trnx.QueryRow(ctx, `
		SELECT ID FROM ROLES
		WHERE ROLE = 'user'
	`).Scan(&roleId)
	if err != nil {
		return &resp, fmt.Errorf("error getting ID from roles : %w", err)
	}

	_, err = trnx.Exec(ctx, `
		INSERT INTO USER_ROLES(USER_ID, ROLE_ID)
		VALUES($1, $2)	
	`, userId, roleId)

	if err != nil {
		return &resp, fmt.Errorf("error inserting data into user_roles : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return &resp, fmt.Errorf("error commiting : %w", err)
	}

	return &resp, nil
}

func (r *Repo) EmailExists(ctx context.Context, email string) (bool, error) {

	var id int

	err := r.pDB.QueryRow(ctx, `
		SELECT ID FROM USERS
		WHERE EMAIL = $1	
	`, email).Scan(&id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("error checking if email %v exists : %w\n", email, err)
	}

	return true, nil
}

func (r *Repo) FetchHashedPass(ctx context.Context, email string) (string, error) {

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

func (r *Repo) FetchUserIdandRoleID(ctx context.Context, email string) (int, int, error) {

	var userID int
	var roleId int

	err := r.pDB.QueryRow(ctx, `
		SELECT ID, ROLE_ID FROM USERS
		INNER JOIN USER_ROLES
		ON USERS.ID = USER_ROLES.USER_ID
		WHERE EMAIL = $1
	`, email).Scan(&userID, &roleId)

	if err != nil {
		return 0, 0, fmt.Errorf("error fetching userID and roleID : %w", err)
	}

	return userID, roleId, nil
}

func (r *Repo) FetchUserName(ctx context.Context, email string) (string, error) {

	var name string

	err := r.pDB.QueryRow(ctx, `

		SELECT NAME FROM USERS
		WHERE EMAIL = $1	
	`, email).Scan(&name)
	if err != nil {
		return name, fmt.Errorf("error getting user name : %w", err)
	}

	return name, nil
}

func (r *Repo) FetchNameUserIdRoleId(ctx context.Context, email string) (int, int, string, error) {

	var userID int
	var roleID int
	var name string

	err := r.pDB.QueryRow(ctx, `
		SELECT ID, NAME, ROLE_ID FROM USERS
		INNER JOIN USER_ROLES
		ON USERS.ID = USER_ROLES.USER_ID
		WHERE EMAIL = $1
	`, email).Scan(&userID, &name, &roleID)

	if err != nil {
		return userID, roleID, name, fmt.Errorf("error getting id, name, role_id : %w", err)
	}

	return userID, roleID, name, nil
}

// refresh -> needs uuid -> cannot get userId
// logout -> needs token -> gets id from token

// logs in -> uuid, refesh and access will be generated
//         -> user_id:id:uuid uuid
//         -> uuid_refresh refresh

// redis

func (r *Repo) RefreshExists(ctx context.Context, userId int) (bool, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)
	// get uuid

	uuid, err := r.rDB.Get(ctx, uuidKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("error checking if uuid exists for user : %v", userId)
	}

	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	_, err = r.rDB.Get(ctx, refreshKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("error getting the refresh token for the user : %v", userId)
	}

	return true, nil
}

func (r *Repo) GetUUID(ctx context.Context, userId int) (string, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	uuid, err := r.rDB.Get(ctx, uuidKey).Result()
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (r *Repo) SetUUID(ctx context.Context, uuid uuid.UUID, userId int) error {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	err := r.rDB.Set(ctx, uuidKey, uuid, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting the uuid : %w", err)
	}

	return nil
}

func (r *Repo) Logout(ctx context.Context, userId int, uuid string) error {
	// del refresh
	refreshKey := fmt.Sprintf("%v_refresh", uuid)
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	pipe := r.rDB.Pipeline()

	pipe.Del(ctx, refreshKey)
	pipe.Del(ctx, uuidKey)
	// del uuid

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting the refresh token after logout : %w\n", err)
	}

	return nil
}

func (r *Repo) SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId int) error {
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)
	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	pipe := r.rDB.Pipeline()

	pipe.Set(ctx, uuidKey, uuid, 0)
	pipe.Set(ctx, refreshKey, Refreshtoken, 0)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error setting refresh and uuid for the user %v : %w\n", userId, err)
	}

	return nil
}

func (r *Repo) DelRefreshToken(ctx context.Context, uuid string) error {

	refreshKey := fmt.Sprintf("%v_refresh", uuid)
	// key := fmt.Sprintf("user:%v:refresh", userId)

	err := r.rDB.Del(ctx, refreshKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting the refersh token : %w", err)
	}

	return nil
}

func (r *Repo) GetRefreshToken(ctx context.Context, uuid uuid.UUID) (string, error) {

	var refreshToken string

	key := fmt.Sprintf("%v_refresh", uuid)

	err := r.rDB.Get(ctx, key).Scan(&refreshToken)
	if err != nil {
		return "", fmt.Errorf("error getting the refresh token : %w\n", err)
	}

	return refreshToken, nil
}
