package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	l "github.com/Zerferrous/Task-Manager-API/internal/logger"
	"github.com/Zerferrous/Task-Manager-API/internal/transport/http/middleware"
)

type Server struct {
	mux        *http.ServeMux
	config     Config
	logger     *l.Logger
	middleware []middleware.Middleware
}

func NewServer(config Config, logger *l.Logger, middlewares ...middleware.Middleware) *Server {
	return &Server{
		mux:        http.NewServeMux(),
		config:     config,
		logger:     logger,
		middleware: middlewares,
	}
}

func (s *Server) RegisterRouters(router ...*Router) {
	for _, router := range router {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (s *Server) Run(ctx context.Context) error {
	mux := middleware.Chain(s.mux, s.middleware...)
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.logger.Info("Sarting HTTP server", slog.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("http server error: %w", err)
		}
	case <-ctx.Done():
		s.logger.Info("Shutting down HTTP server", slog.String("addr", s.config.Addr))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("http server shutdown error: %w", err)
		}

		s.logger.Info("Stopped HTTP server", slog.String("addr", s.config.Addr))
	}

	return nil
}
