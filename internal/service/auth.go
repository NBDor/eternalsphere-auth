package service

import (
	"github.com/NBDor/eternalsphere-auth/internal/models"
	"github.com/NBDor/eternalsphere-auth/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(req *models.RegisterRequest) error
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
}

type AuthService struct {
	userRepo *postgres.UserRepository
}

func NewAuthService(userRepo *postgres.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, err
	}

	// TODO: Generate JWT tokens
	return &models.AuthResponse{
		Token:        "dummy-token",
		RefreshToken: "dummy-refresh-token",
	}, nil
}
