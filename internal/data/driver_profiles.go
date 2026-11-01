package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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

func (m *DriverProfileModel) UpdateAvailability(ctx context.Context, userID string, isAvailable bool) (*DriverProfile, error) {
	stmt := `UPDATE driver_profiles SET is_available = $2 WHERE user_id = $1 RETURNING user_id, vehicle_type, is_available`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{userID, isAvailable}
	var driver DriverProfile
	err := m.DB.QueryRow(ctx, stmt, args...).Scan(&driver.UserID, &driver.VehicleType, &driver.IsAvailable)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &driver, nil
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
