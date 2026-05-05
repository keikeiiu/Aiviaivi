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
	AvatarURL    string
	Bio          string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

func CreateUser(ctx context.Context, db *sql.DB, username string, email string, passwordHash string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
	INSERT INTO users(username, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, username, email, password_hash, avatar_url, bio, role, created_at, updated_at
	`, username, email, passwordHash).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
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
	SELECT id, username, email, password_hash, avatar_url, bio, role, created_at, updated_at
	FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
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
	SELECT id, username, email, password_hash, avatar_url, bio, role, created_at, updated_at
	FROM users WHERE email = $1
	`, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

func UpdateUser(ctx context.Context, db *sql.DB, id string, bio string, avatarURL string) (User, error) {
	var u User
	err := db.QueryRowContext(ctx, `
	UPDATE users SET bio = $2, avatar_url = $3, updated_at = now()
	WHERE id = $1
	RETURNING id, username, email, password_hash, avatar_url, bio, role, created_at, updated_at
	`, id, bio, avatarURL).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}

type UserPublic struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) Public() UserPublic {
	return UserPublic{
		ID:        u.ID,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
		Bio:       u.Bio,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
