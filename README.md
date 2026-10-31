# GoNext

A mini ride-hailing application focused on backend development.

## Stack

- Go, Gin, PostgreSQL, pgxpool, Goose
- Planned: RabbitMQ, WebSocket, Next.js

## Current Features

- Rider registration with input validation and bcrypt password hashing
- Case-insensitive duplicate email protection
- User and driver profile database schemas
- Driver registration with transactional user/profile creation
- Login with HS256 JWT access tokens (30-minute default) and Bearer-protected `/auth/me`

## Run Locally

Requires Go, Docker Compose, and the Goose CLI. From the project root:

```sh
docker compose up -d postgres
export DATABASE_URL='postgres://gonext:gonext_local_password@127.0.0.1:5432/gonext?sslmode=disable'
export JWT_SECRET="$(openssl rand -hex 32)"
goose -env=none -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/api
```

The URL above uses the Compose defaults; adjust it to match your local database.
The API reads process environment variables, not `.env` automatically.
`.env.example` contains only PostgreSQL settings for Docker Compose.
Default API address: `http://localhost:4000`.

Keep `JWT_SECRET` private and stable between restarts; changing it invalidates existing tokens.
The secret must contain at least 32 bytes. Override token lifetime with `-jwt-ttl=30m`.
Send access tokens using `Authorization: Bearer <token>`.

## Unit Tests

```sh
go test ./...
```
