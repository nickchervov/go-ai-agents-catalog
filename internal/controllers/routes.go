package controllers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickchervov/go-ai-agents-catalog/internal/service"
)

func Routes(logger *slog.Logger, svc *service.AgentService, validator TokenValidator) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(LoggerMiddleware(logger))
	r.Use(middleware.Recoverer)

	h := NewAgentsHandler(svc)

	r.Post("/dummyLogin", h.DummyLogin)

	r.Get("/_info", h.Healthz)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(validator))

		r.Get("/categories", h.Categories)
		r.Get("/assistants", h.Assistants)
		r.Get("/assistants/{assistantId}", h.AssistantByID)
		r.Post("/assistants/{assistantId}/run", h.RunAssistant)
		r.Get("/runs/my", h.UserRuns)

		r.Group(func(r chi.Router) {
			r.Use(RequiredRoleMiddleware("admin"))

			r.Post("/categories", h.CreateCategory)
			r.Post("/assistants", h.CreateAssistant)
			r.Put("/assistants/{assistantId}", h.UpdateAssistantByID)
			r.Get("/admin/runs", h.AdminRuns)
		})
	})

	return r
}
