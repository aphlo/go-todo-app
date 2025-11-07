package todo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domaintodo "go-todo-app/internal/domain/todo"
	"go-todo-app/internal/httpx"
	usecasetodo "go-todo-app/internal/usecase/todo"
)

// Service exposes the todo use cases consumed by the HTTP handler.
type Service interface {
	ListTodos(ctx context.Context) ([]domaintodo.Todo, error)
	CreateTodo(ctx context.Context, title string) (domaintodo.Todo, error)
}

// Handler routes todo related HTTP requests.
type Handler struct {
	service Service
}

// NewHandler returns an http.Handler for todo routes.
func NewHandler(service Service) http.Handler {
	return &Handler{service: service}
}

// ServeHTTP muxes GET/POST on /todos.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTodos(w, r)
	case http.MethodPost:
		h.createTodo(w, r)
	default:
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) listTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	todos, err := h.service.ListTodos(ctx)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list todos")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, todos)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var payload struct {
		Title string `json:"title"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.service.CreateTodo(ctx, payload.Title)
	if err != nil {
		if errors.Is(err, usecasetodo.ErrTitleRequired) {
			httpx.WriteError(w, http.StatusBadRequest, "title is required")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create todo")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, todo)
}
