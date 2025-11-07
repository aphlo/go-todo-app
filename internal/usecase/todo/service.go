package todo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	domaintodo "go-todo-app/internal/domain/todo"
)

// Repository abstracts persistence operations for todos.
type Repository interface {
	List(ctx context.Context) ([]domaintodo.Todo, error)
	Create(ctx context.Context, title string) (domaintodo.Todo, error)
}

// ErrTitleRequired indicates validation failure for empty titles.
var ErrTitleRequired = errors.New("title is required")

// Service contains todo-related business logic.
type Service struct {
	repo Repository
}

// NewService wires a repository implementation into the use case layer.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ListTodos fetches all todos.
func (s *Service) ListTodos(ctx context.Context) ([]domaintodo.Todo, error) {
	todos, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	return todos, nil
}

// CreateTodo validates the title and persists a new todo.
func (s *Service) CreateTodo(ctx context.Context, title string) (domaintodo.Todo, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return domaintodo.Todo{}, ErrTitleRequired
	}
	todo, err := s.repo.Create(ctx, trimmed)
	if err != nil {
		return domaintodo.Todo{}, fmt.Errorf("create todo: %w", err)
	}
	return todo, nil
}
