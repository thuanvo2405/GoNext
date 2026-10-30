-- +goose Up
CREATE TABLE driver_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    vehicle_type TEXT NOT NULL CHECK (
        vehicle_type IN ('motorcycle', 'car', 'electric_car')
    ),
    is_available BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE driver_profiles;
