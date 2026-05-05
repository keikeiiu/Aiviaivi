package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

func CreateUser(ctx context.Context, db *sql.DB, username string, email string, passwordHash string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
INSERT INTO users(username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, password_hash, created_at
`, username, email, passwordHash).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrUserExists
		}
		return User{}, err
	}
	return u, nil
}

func GetUserByID(ctx context.Context, db *sql.DB, id string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
SELECT id, username, email, password_hash, created_at
FROM users
WHERE id = $1
`, id).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
SELECT id, username, email, password_hash, created_at
FROM users
WHERE email = $1
`, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}
