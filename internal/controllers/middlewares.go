package controllers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickchervov/go-ai-agents-catalog/internal/domain"
)

type responseWriter struct {
	http.ResponseWriter
	WrittenBytes int
	Status       int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.Status == 0 {
		rw.Status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.WrittenBytes += n
	return n, err
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, Status: http.StatusOK}

			next.ServeHTTP(ww, r)

			duration := time.Since(start).Milliseconds()

			reqId := middleware.GetReqID(r.Context())

			args := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status),
				slog.Int64("duration_ms", duration),
				slog.Int("bytes", ww.WrittenBytes),
			}

			if reqId != "" {
				args = append(args, slog.String("request_id", reqId))
			}

			if ww.Status >= 500 {
				logger.Error("HTTP request failed", args...)
			} else if ww.Status >= 400 {
				logger.Warn("HTTP request warning", args...)
			} else {
				logger.Info("HTTP request processed", args...)
			}
		})
	}
}

func AuthMiddleware(validator TokenValidator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization != "" {
				bearerToken := strings.Split(authorization, " ")
				if len(bearerToken) == 2 {
					claims, err := validator.ValidateToken(bearerToken[1])
					if err == nil {
						ctx := r.Context()
						ctx = context.WithValue(ctx, "userId", claims.UserID)
						ctx = context.WithValue(ctx, "role", claims.Role)

						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
			JSON(w, http.StatusUnauthorized, map[string]ErrorMessage{"error": {Code: domain.ErrCodeUnauthorized, Message: "unauthorized"}})
		})
	}
}

func RequiredRoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role")

			if role != requiredRole {
				JSON(w, http.StatusForbidden, map[string]ErrorMessage{"error": {Code: domain.ErrCodeForbidden, Message: "forbidden"}})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
