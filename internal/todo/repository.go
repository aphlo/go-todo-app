package todo

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository provides DB access for todo entities.
type Repository struct {
	db *sql.DB
}

// NewRepository constructs a Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// List returns all todos ordered by creation time (newest first).
func (r *Repository) List(ctx context.Context) ([]Todo, error) {
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

	var todos []Todo
	for rows.Next() {
		var t Todo
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
func (r *Repository) Create(ctx context.Context, title string) (Todo, error) {
	const query = `
        INSERT INTO todos (title)
        VALUES ($1)
        RETURNING id, title, completed, created_at
    `

	var t Todo
	if err := r.db.QueryRowContext(ctx, query, title).Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
		return Todo{}, fmt.Errorf("insert todo: %w", err)
	}

	return t, nil
}
