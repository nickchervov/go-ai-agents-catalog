package service_test

import "errors"

var (
	ErrInternal    = errors.New("internal error")
	ErrorJWT       = errors.New("jwt error")
	ErrLLMProvider = errors.New("llm provider error")
)
