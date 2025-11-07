package httpserver

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	httpadtodo "go-todo-app/internal/adapter/http/todo"
	"go-todo-app/internal/httpx"
	infratodo "go-todo-app/internal/infrastructure/db/todo"
	usecasetodo "go-todo-app/internal/usecase/todo"
)

// Server wraps the http.Server with app specific handlers and helpers.
type Server struct {
	httpServer *http.Server
}

// New constructs a configured HTTP server instance.
func New(addr string, db *sql.DB) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthHandler(db))

	todoRepo := infratodo.NewPostgresRepository(db)
	todoService := usecasetodo.NewService(todoRepo)
	todoHandler := httpadtodo.NewHandler(todoService)
	mux.Handle("/todos", todoHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{httpServer: srv}
}

// Start begins serving HTTP requests.
func (s *Server) Start() error {
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return err
}

// Shutdown attempts a graceful stop of the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			log.Printf("healthz ping failed: %v", err)
			httpx.WriteError(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}

		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.status = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.status, duration)
	})
}
