package todo

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go-todo-app/internal/httpx"
)

// Handler routes todo related HTTP requests.
type Handler struct {
	repo *Repository
}

// NewHandler returns an http.Handler for todo routes.
func NewHandler(repo *Repository) http.Handler {
	return &Handler{repo: repo}
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

	todos, err := h.repo.List(ctx)
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

	if strings.TrimSpace(payload.Title) == "" {
		httpx.WriteError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo, err := h.repo.Create(ctx, strings.TrimSpace(payload.Title))
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create todo")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, todo)
}
