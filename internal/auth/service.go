package auth

import (
	"context"
	"errors"
	platform_errors "mangahub/pkg/errors"
	"mangahub/pkg/models"
	gormHelper "mangahub/pkg/utils/gorm"
	"time"

	"gorm.io/gorm"
)

type TokenInfo struct {
	AccessToken string
	ExpiresAt   time.Time
	TokenType   string
}

// Defines the contract for accessing user data creation
type UserDataCreator interface {
	CreateUser(ctx context.Context, username string, hashedPassword string) (*models.User, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// Defines the contract for JWT operations
type TokenManager interface {
	GenerateJWT(userID string, userName string) (string, time.Time, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// Defines the contract for password operations
type Hasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hash, password string) bool
}

// Defines the business logic workflow contract
type AuthService interface {
	Register(ctx context.Context, username string, password string) (TokenInfo, error)
	Login(ctx context.Context, username string, password string) (TokenInfo, error)
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

func (s *Service) Register(ctx context.Context, username string, password string) (TokenInfo, error) {

	hashed, err := s.Hasher.HashPassword(password)
	if err != nil {
		return TokenInfo{}, err
	}

	user, err := s.UserCreator.CreateUser(ctx, username, hashed)
	if err != nil {
		if gormHelper.IsUniqueConstraintError(err) {
			return TokenInfo{}, ErrUserConflict
		}

		return TokenInfo{}, err
	}

	token, expiresAt, err := s.TokenMgr.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}, nil
}

func (s *Service) Login(ctx context.Context, username string, password string) (TokenInfo, error) {
	user, err := s.UserCreator.FindUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TokenInfo{}, ErrInvalidCredentials
		}

		return TokenInfo{}, platform_errors.ErrDatabaseOperation
	}

	if !s.Hasher.VerifyPassword(user.PasswordHash, password) {
		return TokenInfo{}, ErrInvalidCredentials
	}

	token, expiresAt, err := s.TokenMgr.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return TokenInfo{}, err
	}

	return TokenInfo{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}, nil
}
