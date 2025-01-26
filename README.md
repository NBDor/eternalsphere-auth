# Eternal Sphere Auth Service

Authentication service for the Eternal Sphere platform.

## Prerequisites
- Go 1.21+
- PostgreSQL
- eternalsphere-shared-go library
- golang-migrate

## Setup

```bash
go mod init github.com/yourusername/eternalsphere-auth
go mod tidy
```

### Database Migrations
Install golang-migrate:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations:
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_users_table

# Apply migrations
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up

# Rollback migrations
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" down
```

## Environment Variables
```bash
# Server
SERVER_ADDRESS=:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
```

## Project Structure
```
eternalsphere-auth/
├── cmd/
│   └── main.go           # Application entry point
├── internal/            # Internal packages
├── migrations/         # Database migrations
└── README.md
```

## API Endpoints

### POST /register
Register new user
```json
{
    "username": "user",
    "email": "user@example.com",
    "password": "password123"
}
```

### POST /login
Authenticate user
```json
{
    "username": "user",
    "password": "password123"
}
```
Response:
```json
{
    "token": "jwt-access-token",
    "refresh_token": "jwt-refresh-token"
}
```

### POST /refresh
Refresh JWT token using refresh token body
Response:
```json
{
    "token": "new-jwt-access-token",
    "refresh_token": "new-jwt-refresh-token"
}
```

## Testing
```bash
go test ./...
```