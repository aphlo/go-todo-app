package todo

import (
	"context"
	"net/http"
	"time"

	domaintodo "go-todo-app/internal/domain/todo"
	"go-todo-app/internal/httpx"
)

// Service exposes the todo use cases consumed by the HTTP handler.
type Service interface {
	ListTodos(ctx context.Context) ([]domaintodo.Todo, error)
}

// Handler routes todo related HTTP requests.
type Handler struct {
	service Service
}

// NewHandler returns an http.Handler for todo routes.
func NewHandler(service Service) http.Handler {
	return &Handler{service: service}
}

// ServeHTTP godoc
// @Summary List todos
// @Tags todos
// @Produce json
// @Success 200 {array} domaintodo.Todo
// @Failure 405 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /todos [get]
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	todos, err := h.service.ListTodos(ctx)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list todos")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, todos)
}
