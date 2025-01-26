# Eternal Sphere Auth Service

Authentication service for the Eternal Sphere platform.

## Prerequisites
- Go 1.21+
- PostgreSQL
- eternalsphere-shared-go library

## Setup
```bash
go mod init github.com/yourusername/eternalsphere-auth
go mod tidy
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

### POST /refresh
Refresh JWT token using refresh token

## Testing
```bash
go test ./...
```