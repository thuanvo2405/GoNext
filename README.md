# GoNext

A mini ride-hailing application focused on backend development.

## Stack

- Go, Gin, PostgreSQL, pgxpool, Goose
- Planned: RabbitMQ, WebSocket, Next.js

## Current Features

- Rider registration with input validation and bcrypt password hashing
- Case-insensitive duplicate email protection
- User and driver profile database schemas

## Run Locally

Requires Go, Docker Compose, and the Goose CLI. From the project root:

```sh
docker compose up -d postgres
export DATABASE_URL='postgres://gonext:gonext_local_password@127.0.0.1:5432/gonext?sslmode=disable'
goose -env=none -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/api
```

The URL above uses the Compose defaults; adjust it to match your local database.
The API reads process environment variables, not `.env` automatically.
Default API address: `http://localhost:4000`.

## Unit Tests

```sh
go test ./...
```
