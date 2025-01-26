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

## Project Structure
```
eternalsphere-auth/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── handlers/           # HTTP request handlers
│   ├── models/             # Data models
│   ├── repository/         # Database operations
│   └── service/            # Business logic
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