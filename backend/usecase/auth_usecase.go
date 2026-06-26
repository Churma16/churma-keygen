package usecase

import (
	"errors"
	"time"

	"churma-keygen/backend/domain"
	"churma-keygen/backend/dtos"
	"churma-keygen/backend/middleware"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(req dtos.LoginRequest) (*dtos.LoginResponse, error)
	GetMe(userID, username, role string) (*dtos.UserResponse, error)
	UpdateProfile(userID string, req dtos.UpdateProfileRequest) error
}

type authUsecaseImpl struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(userRepo domain.UserRepository) AuthUsecase {
	return &authUsecaseImpl{userRepo: userRepo}
}

func (s *authUsecaseImpl) Login(req dtos.LoginRequest) (*dtos.LoginResponse, error) {
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

func (s *authUsecaseImpl) GetMe(userID, username, role string) (*dtos.UserResponse, error) {
	return &dtos.UserResponse{
		ID:       userID,
		Username: username,
		Role:     role,
	}, nil
}

func (s *authUsecaseImpl) UpdateProfile(userID string, req dtos.UpdateProfileRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		return errors.New("password saat ini tidak valid")
	}

	if req.Username != user.Username {
		count, err := s.userRepo.CountByUsername(req.Username)
		if err != nil {
			return errors.New("gagal memvalidasi username")
		}
		if count > 0 {
			return errors.New("username sudah digunakan oleh akun lain")
		}
		user.Username = req.Username
	}

	if req.NewPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("gagal memproses password baru")
		}
		user.PasswordHash = string(hashedPassword)
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return errors.New("gagal memperbarui profil")
	}

	return nil
}
