package mocks

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type MockJWT struct {
	Err    error
	TTL    time.Duration
	Secret []byte
}

func (j *MockJWT) CreateJWTToken(userID, role, email string) (string, error) {
	if j.Err != nil {
		return "", j.Err
	}

	claims := domain.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TTL)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(j.Secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}
