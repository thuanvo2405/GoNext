package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DriverProfile struct {
	UserID      string `json:"user_id"`
	VehicleType string `json:"vehicle_type"`
	IsAvailable bool   `json:"is_available"`
}

type DriverProfileModel struct {
	DB *pgxpool.Pool
}

func (m *DriverProfileModel) Insert(ctx context.Context, user *User, driver *DriverProfile) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userID string

	stmt := `INSERT INTO users(email, password_hash, role) VALUES($1, $2, $3) RETURNING id`

	args := []any{user.Email, user.Password, user.Role}

	err = tx.QueryRow(ctx, stmt, args...).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "users_email_unique" {
			return ErrDuplicateEmail
		}

		return err
	}

	stmt = `INSERT INTO driver_profiles(user_id, vehicle_type) VALUES($1, $2)`

	args = []any{userID, driver.VehicleType}

	_, err = tx.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	user.ID = userID
	driver.UserID = userID
	driver.IsAvailable = false

	return nil
}
