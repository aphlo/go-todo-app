package todo

import (
	"context"
	"fmt"

	domaintodo "go-todo-app/internal/domain/todo"
)

// Repository abstracts persistence operations for todos.
type Repository interface {
	List(ctx context.Context) ([]domaintodo.Todo, error)
}

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
