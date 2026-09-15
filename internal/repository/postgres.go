package repository

import (
	"context"
	"fmt"

	"github.com/ImmortaL-jsdev/task-manager/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	var created models.Task

	query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id, title, description, status, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status).Scan(&created.ID, &created.Title, &created.Description, &created.Status, &created.CreatedAt, &created.UpdatedAt)

	if err != nil {
		return models.Task{}, fmt.Errorf("failed to insert task: %w", err)
	}

	return created, nil
}

func (s *PostgresStore) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}

func (s *PostgresStore) GetTaskByID(ctx context.Context, id string) (models.Task, error) {
	var task models.Task
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1`

	err := s.pool.QueryRow(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to get task by id: %w", err)
	}
	return task, nil
}
