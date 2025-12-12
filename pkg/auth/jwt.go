package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"Dash/pkg/env"
	"Dash/pkg/logger"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewBaseClaim(userID string) *Claims {
	return &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
}

func newClaim(baseClaim *Claims, expirationTime time.Time) *Claims {
	return &Claims{
		UserID: baseClaim.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  baseClaim.IssuedAt,
		},
	}
}

func GenerateAccessToken(baseClaim *Claims) (string, error) {
	secret := getSecret()
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := newClaim(baseClaim, expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Error("GenerateAccessToken:SIGN_FAILED", map[string]interface{}{
			"user_id": baseClaim.UserID,
			"error":   err.Error(),
		})
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(baseClaim *Claims) (string, error) {
	secret := getSecret()
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := newClaim(baseClaim, expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Error("GenerateRefreshToken:SIGN_FAILED", map[string]interface{}{
			"user_id": baseClaim.UserID,
			"error":   err.Error(),
		})
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	secret := getSecret()
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.Warn("ValidateToken:TOKEN_EXPIRED", map[string]interface{}{
				"error": err.Error(),
			})
			return "", ErrTokenExpired
		}
		logger.Error("ValidateToken:TOKEN_INVALID", map[string]interface{}{
			"error": err.Error(),
		})
		return "", ErrInvalidToken
	}

	if !token.Valid {
		logger.Error("ValidateToken:TOKEN_NOT_VALID", map[string]interface{}{})
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}

func getSecret() string {
	secret := env.GetString("JWT_SECRET")
	if secret == "" {
		logger.Warn("getSecret:DEFAULT_SECRET_USED", map[string]interface{}{
			"warning": "JWT_SECRET not set, using default secret (change in production)",
		})
		return "4a8f3b9c2d8e4f1a6b5c9d2e7f3a8b4c1d6e9f2a5b8c3d7e1f4a9b2c6d8e3f7a1b4c9d2e6f"
	}
	return secret
}
