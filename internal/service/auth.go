package service

import (
	"errors"

	"github.com/NBDor/eternalsphere-auth/internal/models"
	"github.com/NBDor/eternalsphere-auth/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(req *models.RegisterRequest) error
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(token string) (*models.AuthResponse, error)
}

type AuthService struct {
	userRepo  *postgres.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *postgres.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
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
		return nil, errors.New("invalid credentials")
	}

	tokenPair, err := generateTokenPair(user.ID, user.Username, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(token string) (*models.AuthResponse, error) {
	userID, username, err := validateRefreshToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	tokenPair, err := generateTokenPair(userID, username, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
