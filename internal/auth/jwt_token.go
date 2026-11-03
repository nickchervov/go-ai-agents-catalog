package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type JWTProvider struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTProvider(secret []byte, ttl time.Duration) *JWTProvider {
	return &JWTProvider{secret: secret, ttl: ttl}
}

func (j *JWTProvider) CreateJWTToken(userID, role, email string) (string, error) {
	claims := domain.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (j *JWTProvider) ValidateToken(token string) (*domain.Claims, error) {
	claims := domain.Claims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signed method")
		}
		return j.secret, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, fmt.Errorf("validating token: %w", err)
	}
	return &claims, nil
}
