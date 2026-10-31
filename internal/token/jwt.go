package token

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) (*Manager, error) {
	if len(secret) < 32 || strings.TrimSpace(secret) == "" {
		return nil, errors.New("JWT secret must contain at least 32 bytes; use a cryptographically random secret")
	}
	if ttl <= 0 {
		return nil, errors.New("JWT lifetime must be greater than zero")
	}
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
	}, nil
}

func (m *Manager) Sign(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(m.ttl).Unix(),
	})

	return token.SignedString(m.secret)
}

func (m *Manager) Verify(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return m.secret, nil
	})
	if err != nil {
		return "", errors.New("Could not parse token.")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return "", errors.New("Invalid token!")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Invalid token claims.")
	}

	// email := claims["email"].(string)
	userId, ok := claims["userId"].(string)
	if !ok || strings.TrimSpace(userId) == "" {
		return "", errors.New("invalid token claims")
	}
	return userId, nil
}
