package service_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
	"github.com/nickchervov/go-ai-agents-catalog/internal/dto"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
	"github.com/nickchervov/go-ai-agents-catalog/mocks"
)

func TestDummyLogin(t *testing.T) {
	adminID, err := uuid.Parse(domain.RoleDatas["admin"].UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	userID, err := uuid.Parse(domain.RoleDatas["user"].UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name           string
		repoError      error
		jwtError       error
		input          dto.DummyLoginInput
		expectedOutput *dto.DummyLoginOutput
		checkErr       func(error) bool
	}{
		{
			name:      "good case admin",
			repoError: nil,
			jwtError:  nil,
			input: dto.DummyLoginInput{
				Role: "admin",
			},
			expectedOutput: &dto.DummyLoginOutput{
				User: struct {
					ID    uuid.UUID "json:\"id\""
					Email string    "json:\"email\""
					Role  string    "json:\"role\""
				}{
					ID:    adminID,
					Email: "admin@example.com",
					Role:  "admin",
				},
			},
			checkErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name:      "good case user",
			repoError: nil,
			jwtError:  nil,
			input: dto.DummyLoginInput{
				Role: "user",
			},
			expectedOutput: &dto.DummyLoginOutput{
				User: struct {
					ID    uuid.UUID "json:\"id\""
					Email string    "json:\"email\""
					Role  string    "json:\"role\""
				}{
					ID:    userID,
					Email: "user@example.com",
					Role:  "user",
				},
			},
			checkErr: func(err error) bool {
				return err == nil
			},
		},
		{
			name:      "incorrect role",
			repoError: nil,
			jwtError:  nil,
			input: dto.DummyLoginInput{
				Role: "adadad",
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, domain.ErrIncorrectRole)
			},
		},
		{
			name:      "jwt error",
			repoError: nil,
			jwtError:  ErrorJWT,
			input: dto.DummyLoginInput{
				Role: "admin",
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrorJWT)
			},
		},
		{
			name:      "internal error",
			repoError: ErrInternal,
			jwtError:  nil,
			input: dto.DummyLoginInput{
				Role: "admin",
			},
			expectedOutput: nil,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrInternal)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockRepo{
				Err: tt.repoError,
			}

			mockLLM := &mocks.MockLLM{
				Result: "Nice prompt take result back",
			}

			mockJWT := &mocks.MockJWT{
				Err:    tt.jwtError,
				Secret: []byte("qwerty"),
				TTL:    10 * time.Minute,
			}

			svc := service.NewAgentService(mockRepo, mockJWT, mockLLM)

			output, err := svc.DummyLogin(t.Context(), tt.input)
			if err != nil {
				if !tt.checkErr(err) {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if output.Token == "" {
					t.Fatal("token is empty")
				}

				var claims domain.Claims
				_, err := jwt.ParseWithClaims(output.Token, &claims, func(t *jwt.Token) (any, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("incorrect signed method")
					}
					return mockJWT.Secret, nil
				})

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if claims.UserID != tt.expectedOutput.User.ID.String() {
					t.Fatalf("got user id: %s want: %v", claims.UserID, tt.expectedOutput.User.ID)
				}

				if claims.Role != tt.expectedOutput.User.Role {
					t.Fatalf("got user role: %s want: %s", claims.Role, tt.expectedOutput.User.Role)
				}

				if claims.Email != tt.expectedOutput.User.Email {
					t.Fatalf("got user email: %s want: %v", claims.Email, tt.expectedOutput.User.Email)
				}
			}
		})
	}
}
