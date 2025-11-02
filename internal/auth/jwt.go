package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	SecretKey string
	Duration  time.Duration
}

// Constructor, injecting the secret and duration
func NewJWTManager(secret string, durationStr string) (*JWTManager, error) {
	if len(secret) == 0 {
		return nil, errors.New("JWT_SECRET environment variable is not set or empty")
	}

	duration := defaultTokenDuration
	if durationStr != "" {
		d, err := time.ParseDuration(durationStr)
		if err != nil {
			return nil, errors.New("invalid JWT_ACCESS_TOKEN_LIFETIME format: " + err.Error())
		}
		duration = d
	}

	return &JWTManager{
		SecretKey: secret,
		Duration:  duration,
	}, nil
}

const defaultTokenDuration = 8 * time.Hour

type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func (m *JWTManager) GenerateJWT(userID, userName string) (string, int64, error) {
	claims := JWTClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "mangahub-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.SecretKey))
	if err != nil {
		return "", 0, ErrTokenGeneration
	}

	return signed, int64(m.Duration.Seconds()), nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
