package service

import (
	"context"
	"fmt"

	myerrors "github.com/ImmortaL-jsdev/task-manager/internal/errors"
	"github.com/ImmortaL-jsdev/task-manager/internal/models"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id string) (models.Task, error)
	UpdateTask(ctx context.Context, id string, task models.Task) (models.Task, error)
	DeleteTask(ctx context.Context, id string) error
}
type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	if task.Title == "" {
		return models.Task{}, &myerrors.ValidationError{Message: "title cannot be empty"}
	}
	if task.Status == "" {
		task.Status = "todo"
	}
	switch task.Status {
	case "todo", "in_progress", "done":

	default:
		return models.Task{}, &myerrors.ValidationError{Message: "invalid status"}
	}

	created, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to create task: %w", err)
	}
	return created, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAllTasks(ctx)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (models.Task, error) {
	return s.repo.GetTaskByID(ctx, id)
}
func (s *TaskService) UpdateTask(ctx context.Context, id string, task models.Task) (models.Task, error) {
	if task.Title == "" {
		return models.Task{}, &myerrors.ValidationError{Message: "title cannot be empty"}
	}
	if task.Status == "" {
		task.Status = "todo"
	}
	switch task.Status {
	case "todo", "in_progress", "done":

	default:
		return models.Task{}, &myerrors.ValidationError{Message: "invalid status"}
	}
	updated, err := s.repo.UpdateTask(ctx, id, task)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to update task: %w", err)
	}
	return updated, nil
}
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {

	return s.repo.DeleteTask(ctx, id)
}
