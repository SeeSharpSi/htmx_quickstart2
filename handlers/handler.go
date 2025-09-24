package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"seesharpsi/web_roguelike/logger"
	"seesharpsi/web_roguelike/session"
	"seesharpsi/web_roguelike/templ"
)

type Handler struct {
	Manager *session.Manager
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Info("handling index request", "request_id", requestID)

	// replace _ with sess if you want to use a session
	_, cookie := h.Manager.GetOrCreateSession(r)
	http.SetCookie(w, &cookie)

	if err := templ.Index().Render(context.Background(), w); err != nil {
		slog.Error("failed to render index template", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Info("handling test request", "request_id", requestID)

	// replace _ with sess if you want to use a session
	_, cookie := h.Manager.GetOrCreateSession(r)
	http.SetCookie(w, &cookie)

	if err := templ.Test().Render(context.Background(), w); err != nil {
		slog.Error("failed to render test template", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Info("handling 404 request", "request_id", requestID, "path", r.URL.Path)

	w.WriteHeader(http.StatusNotFound)
	if err := templ.Error404().Render(context.Background(), w); err != nil {
		slog.Error("failed to render 404 template", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Debug("handling health check", "request_id", requestID)

	// Return simple JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}
