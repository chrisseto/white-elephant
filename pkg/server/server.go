package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	// "github.com/go-chi/httptracer"
)

type Server struct {
	http.Server

	db *sqlx.DB
}

func New(db *sqlx.DB) *Server {
	api := &API{db: db}

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(httptracer.Tracer())
	r.Use(middleware.AllowContentEncoding(
		"application/json",
	))

	r.Post("/api/rooms/", api.newRoom)
	r.Get("/api/rooms/{id}/", api.getRoom)

	// Serve our static assets
	r.Get("/static/{path}", assets)

	// Catch all for the FE
	r.Get("/", indexHTML)
	r.Get("/*", indexHTML)

	return &Server{
		db: db,
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
