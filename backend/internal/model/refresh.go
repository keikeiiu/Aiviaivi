package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
)

// StoreRefreshToken saves a hashed refresh token for a user.
func StoreRefreshToken(ctx context.Context, db *sql.DB, userID string, rawToken string, expiresAt time.Time) error {
	hash := hashToken(rawToken)
	_, err := db.ExecContext(ctx, `
	INSERT INTO refresh_tokens(user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	`, userID, hash, expiresAt)
	return err
}

// ConsumeRefreshToken validates and removes a refresh token.
// Returns the user ID if valid, or an error if invalid/expired/revoked.
func ConsumeRefreshToken(ctx context.Context, db *sql.DB, rawToken string) (string, error) {
	hash := hashToken(rawToken)

	// Use explicit transaction to guarantee commit
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var userID string
	var expiresAt time.Time
	var revokedAt sql.NullTime

	err = tx.QueryRowContext(ctx, `
	DELETE FROM refresh_tokens
	WHERE token_hash = $1
	RETURNING user_id, expires_at, revoked_at
	`, hash).Scan(&userID, &expiresAt, &revokedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrRefreshTokenNotFound
		}
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	if revokedAt.Valid {
		return "", ErrRefreshTokenRevoked
	}

	if time.Now().After(expiresAt) {
		return "", ErrRefreshTokenExpired
	}

	return userID, nil
}

// RevokeUserRefreshTokens revokes all refresh tokens for a user.
func RevokeUserRefreshTokens(ctx context.Context, db *sql.DB, userID string) error {
	_, err := db.ExecContext(ctx, `
	UPDATE refresh_tokens SET revoked_at = now()
	WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > now()
	`, userID)
	return err
}

// CleanupExpiredTokens removes expired refresh tokens older than the given duration.
func CleanupExpiredTokens(ctx context.Context, db *sql.DB, olderThan time.Duration) (int64, error) {
	res, err := db.ExecContext(ctx, `
	DELETE FROM refresh_tokens
	WHERE expires_at < now() - $1
	`, olderThan)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
