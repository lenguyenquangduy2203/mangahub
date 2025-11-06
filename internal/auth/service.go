package auth

import (
	"context"
	"errors"
	"mangahub/data/models"
	platform_errors "mangahub/pkg/errors"
	"mangahub/pkg/models/dtos"
	gormHelper "mangahub/pkg/utils/gorm"

	"gorm.io/gorm"
)

// Defines the contract for accessing user data creation
type UserDataCreator interface {
	CreateUser(ctx context.Context, username string, hashedPassword string) (*models.User, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// Defines the contract for JWT operations
type TokenManager interface {
	GenerateJWT(userID string, userName string) (string, int64, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// Defines the contract for password operations
type Hasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hash, password string) bool
}

// Defines the business logic workflow contract
type AuthService interface {
	Register(ctx context.Context, username string, password string) (dtos.TokenResponse, error)
	Login(ctx context.Context, username string, password string) (dtos.TokenResponse, error)
}

type Service struct {
	UserCreator UserDataCreator
	TokenMgr    TokenManager
	Hasher      Hasher
}

// Constructor for the AuthService
func NewService(uc UserDataCreator, tm TokenManager, h Hasher) *Service {
	return &Service{
		UserCreator: uc,
		TokenMgr:    tm,
		Hasher:      h,
	}
}

func (s *Service) Register(ctx context.Context, username string, password string) (dtos.TokenResponse, error) {

	hashed, err := s.Hasher.HashPassword(password)
	if err != nil {
		return dtos.TokenResponse{}, err
	}

	user, err := s.UserCreator.CreateUser(ctx, username, hashed)
	if err != nil {
		if gormHelper.IsUniqueConstraintError(err) {
			return dtos.TokenResponse{}, ErrUserConflict
		}

		return dtos.TokenResponse{}, err
	}

	token, expiresIn, err := s.TokenMgr.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return dtos.TokenResponse{}, err
	}

	return dtos.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *Service) Login(ctx context.Context, username string, password string) (dtos.TokenResponse, error) {
	user, err := s.UserCreator.FindUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dtos.TokenResponse{}, ErrInvalidCredentials
		}

		return dtos.TokenResponse{}, platform_errors.ErrDatabaseOperation
	}

	if !s.Hasher.VerifyPassword(user.PasswordHash, password) {
		return dtos.TokenResponse{}, ErrInvalidCredentials
	}

	token, expiresIn, err := s.TokenMgr.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return dtos.TokenResponse{}, err
	}

	return dtos.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}, nil
}
