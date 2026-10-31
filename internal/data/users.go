package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRow(ctx, `
		SELECT id, email, password_hash, role
		FROM users
		WHERE lower(email) = lower($1)
	`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Insert(ctx context.Context, user *User) error {
	stmt := `INSERT INTO users(email, password_hash, role) VALUES($1, $2, $3) RETURNING id, email, role`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{user.Email, user.Password, user.Role}

	err := m.DB.QueryRow(ctx, stmt, args...).Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "users_email_unique" {
			return ErrDuplicateEmail
		}

		return err
	}
	return nil
}
