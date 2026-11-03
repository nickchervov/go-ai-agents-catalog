package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	DSN          string
	LLMAPI       string
	WantMockLLM  bool
	LoggerFormat string
	Secret       []byte
	TTLToken     time.Duration
	BaseURL      string
}

func LoadConfig() (Config, error) {
	var cfg Config

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = ":8080"
	}
	cfg.Port = port

	dsn, ok := os.LookupEnv("DSN")
	if !ok {
		return Config{}, fmt.Errorf("loading config: not found dsn")
	}
	cfg.DSN = dsn

	llm, ok := os.LookupEnv("LLMAPI")
	if !ok || llm == "" {
		cfg.WantMockLLM = true
	}
	cfg.LLMAPI = llm

	format, ok := os.LookupEnv("LOGGER_FORMAT")
	if !ok {
		format = "json"
	}
	cfg.LoggerFormat = format

	secret, ok := os.LookupEnv("SECRET")
	if !ok {
		secret = "sadasdQW"
	}
	cfg.Secret = []byte(secret)

	ttl, ok := os.LookupEnv("TTL_TOKEN_MIN")
	if !ok {
		ttl = "60"
	}
	t, err := strconv.Atoi(ttl)
	if err != nil {
		return Config{}, fmt.Errorf("converting ttl to int: %w", err)
	}
	cfg.TTLToken = time.Duration(t) * time.Minute

	baseURL, ok := os.LookupEnv("BASE_URL")
	if !ok || baseURL == "" {
		cfg.WantMockLLM = true
	}
	cfg.BaseURL = baseURL

	return cfg, nil
}
