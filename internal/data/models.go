package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Users         UserModel
	DriverProfile DriverProfileModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users:         UserModel{DB: db},
		DriverProfile: DriverProfileModel{DB: db},
	}
}
