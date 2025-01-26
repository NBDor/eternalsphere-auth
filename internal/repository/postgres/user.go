package postgres

import (
	"database/sql"
	"errors"

	"github.com/NBDor/eternalsphere-auth/internal/models"
	shared "github.com/NBDor/eternalsphere-shared-go/database/postgres"
)

type UserRepository struct {
	conn *shared.Connection
}

func NewUserRepository(conn *shared.Connection) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, password_hash, email)
		VALUES ($1, $2, $3)
		RETURNING id`

	return r.conn.Transaction(func(tx *sql.Tx) error {
		return tx.QueryRow(query, user.Username, user.PasswordHash, user.Email).Scan(&user.ID)
	})
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, email
		FROM users
		WHERE username = $1`

	user := &models.User{}
	err := r.conn.DB().QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}
