package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/nickchervov/go-ai-agents-catalog/internal/adapters/postgres"
	"github.com/nickchervov/go-ai-agents-catalog/internal/auth"
	"github.com/nickchervov/go-ai-agents-catalog/internal/config"
	"github.com/nickchervov/go-ai-agents-catalog/internal/controllers"
	"github.com/nickchervov/go-ai-agents-catalog/internal/llm"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
	"github.com/nickchervov/go-ai-agents-catalog/mocks"
	"github.com/nickchervov/go-ai-agents-catalog/pkg/logger"
)

func RunApp(ctx context.Context, cfg config.Config) error {
	logger := logger.NewLogger(cfg.LoggerFormat)
	slog.SetDefault(logger)

	repo, err := postgres.NewPostgresRepository(ctx, cfg.DSN)
	if err != nil {
		return fmt.Errorf("creating postgres repo: %w", err)
	}
	defer repo.Close()

	token := auth.NewJWTProvider(cfg.Secret, cfg.TTLToken)

	var llmAgent service.LLM
	if cfg.WantMockLLM {
		llmAgent = &mocks.MockLLM{Result: "Nice prompt take result back"}
	} else {
		llmAgent = llm.NewLLMAgent(cfg.LLMAPI, cfg.BaseURL)
	}

	svc := service.NewAgentService(repo, token, llmAgent)

	r := controllers.Routes(logger, svc, token)

	server := http.Server{
		Addr:         cfg.Port,
		Handler:      r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  35 * time.Second,
		IdleTimeout:  35 * time.Second,
	}

	gsCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		log.Printf("Starting server on port: %s...\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
			return
		}
		log.Println("Server closed.")
	}()

	select {
	case e := <-errChan:
		return fmt.Errorf("listening server: %w", e)
	case <-gsCtx.Done():
		log.Println("Starting graceful shutdown...")

		timeoutCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()

		if err := server.Shutdown(timeoutCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
		return nil
	}
}

func main() {
	_ = godotenv.Load(".env")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	if err := RunApp(context.Background(), cfg); err != nil {
		log.Fatalf("running app: %v", err)
	}
}
