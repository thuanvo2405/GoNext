# GoNext

A mini ride-hailing application built with Go and Gin.

## Local development

Requirements: Go, Docker, and Docker Compose.

### PostgreSQL

Start PostgreSQL:

```sh
docker compose up -d postgres
docker compose ps
```

The default local connection settings are:

- Host: `127.0.0.1`
- Port: `5432`
- Database/user: `gonext`
- Password: `gonext_local_password` (local development only)

Compose works without a `.env` file. To change these defaults, copy
`.env.example` to `.env` and edit it. If port 5432 is already occupied, set
`POSTGRES_PORT=5433`; use that port when connecting from the host.

The container running does not guarantee PostgreSQL is ready. Verify a database
connection inside the container (retry after a few seconds if it is still starting):

```sh
docker compose exec postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT current_database(), current_user;"'
```

Stop the service:

```sh
docker compose down
```

The named volume preserves database data across container recreation.
Do not add `--volumes` to the shutdown command unless you intend to delete
the database data. Initialization settings apply only to an empty data volume;
changing the configured password does not update an existing database user.

### API

Run Go directly on the host:

```sh
go run ./cmd/api
```

In another terminal:

```sh
curl -i http://localhost:8080/health
```

The API currently serves `/health`; database integration is not implemented yet.

### Tests

```sh
go test -v ./...
```
