package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"seesharpsi/web_roguelike/handlers"
	"seesharpsi/web_roguelike/logger"
	"seesharpsi/web_roguelike/session"
)

// custom404Handler wraps a mux to provide custom 404 pages
func custom404Handler(mux *http.ServeMux, h *handlers.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the route exists in the mux
		_, pattern := mux.Handler(r)
		if pattern == "" {
			// Route not found, serve custom 404
			h.NotFound(w, r)
			return
		}
		// Route exists, serve normally
		mux.ServeHTTP(w, r)
	})
}

func main() {
	port := flag.Int("port", 9779, "port the server runs on")
	address := flag.String("address", "http://localhost", "address the server runs on")
	flag.Parse()

	// Setup structured logging
	logger.SetupLogger()

	// ip parsing
	base_ip := *address
	ip := base_ip + ":" + strconv.Itoa(*port)
	root_ip, err := url.Parse(ip)
	if err != nil {
		slog.Error("failed to parse server address", "error", err, "address", ip)
		os.Exit(1)
	}

	sessionManager := session.NewManager()

	h := &handlers.Handler{
		Manager: sessionManager,
	}

	// set up mux with middleware
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/test", h.Test)

	// Custom 404 handler for unmatched routes
	mux.HandleFunc("/404", h.NotFound)

	// Wrap with middleware and custom 404 handler
	handler := logger.PanicRecovery(logger.RequestLogger(custom404Handler(mux, h)))

	server := http.Server{
		Addr:    root_ip.Host,
		Handler: handler,
	}

	// Channel to listen for interrupt signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start server
	slog.Info("starting server", "address", root_ip.Host)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-done
	slog.Info("shutting down server")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
