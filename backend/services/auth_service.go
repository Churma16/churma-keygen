package services

import (
	"errors"
	"time"

	"churma-keygen/backend/dtos"
	"churma-keygen/backend/middleware"
	"churma-keygen/backend/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req dtos.LoginRequest) (*dtos.LoginResponse, error)
	GetMe(userID, username, role string) (*dtos.UserResponse, error)
}

type authServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authServiceImpl{userRepo: userRepo}
}

func (s *authServiceImpl) Login(req dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate session token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.AdminClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "churma-keygen-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		return nil, errors.New("failed to generate session token")
	}

	return &dtos.LoginResponse{
		Token:    tokenString,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

func (s *authServiceImpl) GetMe(userID, username, role string) (*dtos.UserResponse, error) {
	return &dtos.UserResponse{
		ID:       userID,
		Username: username,
		Role:     role,
	}, nil
}
