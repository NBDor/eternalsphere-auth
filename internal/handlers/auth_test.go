package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NBDor/eternalsphere-auth/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *models.RegisterRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful Registration", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewAuthHandler(mockService)

		// Setup test request
		req := models.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		mockService.On("Register", &req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewAuthHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidReq := `{"username": ""}`
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBufferString(invalidReq))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewAuthHandler(mockService)

		req := models.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		mockService.On("Register", &req).Return(errors.New("service error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful Login", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewAuthHandler(mockService)

		req := models.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		expectedResp := &models.AuthResponse{
			Token:        "test-token",
			RefreshToken: "test-refresh",
		}
		mockService.On("Login", &req).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResp.Token, response.Token)
		assert.Equal(t, expectedResp.RefreshToken, response.RefreshToken)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockService := new(MockAuthService)
		handler := NewAuthHandler(mockService)

		req := models.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}
		mockService.On("Login", &req).Return(&models.AuthResponse{}, errors.New("invalid credentials"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})
}
