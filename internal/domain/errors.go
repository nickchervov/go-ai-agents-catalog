package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrIncorrectID          = errors.New("incorrect id")
	ErrIncorrectUserID      = errors.New("incorrect user id")
	ErrIncorrectAssistantID = errors.New("incorrect assistant id")
	ErrAssistantNotFound    = errors.New("assistant not found")
	ErrIncorrectRole        = errors.New("incorrect role")
	ErrAssistantInactive    = errors.New("assistant inactive")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrEmptyUserPrompt      = errors.New("user prompt must be not empty")
	ErrLLMProvider          = errors.New("LLM provider issue")
	ErrEmptyIsActivated     = errors.New("is activated must be not empty")
)

var (
	ErrCodeInvalidRequest    = "INVALID_REQUEST"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeCategoryNotFound  = "CATEGORY_NOT_FOUND"
	ErrCodeAssistantNotFound = "ASSISTANT_NOT_FOUND"
	ErrCodeAssistantInactive = "ASSISTANT_INACTIVE"
	ErrCodeLLMProvider       = "LLM_PROVIDER_ERROR"
	ErrCodeInternal          = "INTERNAL_ERROR"
)

type ValidationErr struct {
	Fields map[string]string
}

func (e *ValidationErr) Error() string {
	var errs []string
	for field, msg := range e.Fields {
		errs = append(errs, fmt.Sprintf("%s:%s", field, msg))
	}
	return "validation failed: " + strings.Join(errs, "; ")
}
