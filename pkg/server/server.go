package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	http.Server

	router chi.Router
}

func New() *Server {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", indexHTML)
	r.Get("/static/{path}", assets)

	return &Server{
		Server: http.Server{
			Handler: r,
		},
	}
}

func (s *Server) Serve(ctx context.Context, lis net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		_ = s.Server.Shutdown(shutdownCtx)
	}()

	err := s.Server.Serve(lis)

	// ErrServerClosed indicates the server shutdown gracefully.
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return errors.Wrap(err, "serving HTTP")
}
