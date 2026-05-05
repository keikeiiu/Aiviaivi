package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(userID string, secret string, expires time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(expires).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func NewRefreshToken(userID string, secret string, expires time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  now.Unix(),
		"exp":  now.Add(expires).Unix(),
		"type": "refresh",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func ParseToken(tokenStr string, secret string) (*jwt.Token, jwt.MapClaims, error) {
	var claims jwt.MapClaims
	t, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, nil, err
	}
	return t, claims, nil
}
