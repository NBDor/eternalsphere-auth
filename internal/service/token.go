package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func generateTokenPair(userID int64, username string, secret string) (*TokenPair, error) {
	// Access token - short lived (15 minutes)
	accessToken, err := generateToken(userID, username, secret, time.Minute*15)
	if err != nil {
		return nil, err
	}

	// Refresh token - longer lived (7 days)
	refreshToken, err := generateToken(userID, username, secret, time.Hour*24*7)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateToken(userID int64, username, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateRefreshToken(tokenString, secret string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", jwt.ErrInvalidKey
	}

	userID := int64(claims["user_id"].(float64))
	username := claims["username"].(string)

	return userID, username, nil
}
