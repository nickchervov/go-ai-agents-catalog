package domain

import "github.com/golang-jwt/jwt/v5"

var RoleDatas = map[string]struct {
	UserID string
	Email  string
}{
	"admin": {"11111111-1111-1111-1111-111111111111", "admin@example.com"},
	"user":  {"22222222-2222-2222-2222-222222222222", "user@example.com"},
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
