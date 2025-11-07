package todo

import (
	"context"
	"database/sql"
	"fmt"

	domaintodo "go-todo-app/internal/domain/todo"
)

// PostgresRepository persists todos in a Postgres database.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository returns a Repository backed by Postgres.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// List returns all todos ordered by creation time (newest first).
func (r *PostgresRepository) List(ctx context.Context) ([]domaintodo.Todo, error) {
	const query = `
        SELECT id, title, completed, created_at
        FROM todos
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	var todos []domaintodo.Todo
	for rows.Next() {
		var t domaintodo.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todos: %w", err)
	}

	return todos, nil
}

// Create inserts a new todo with the provided title.
func (r *PostgresRepository) Create(ctx context.Context, title string) (domaintodo.Todo, error) {
	const query = `
        INSERT INTO todos (title)
        VALUES ($1)
        RETURNING id, title, completed, created_at
    `

	var t domaintodo.Todo
	if err := r.db.QueryRowContext(ctx, query, title).Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
		return domaintodo.Todo{}, fmt.Errorf("insert todo: %w", err)
	}

	return t, nil
}
