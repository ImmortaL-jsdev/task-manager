package service_test

import (
	"context"
	"errors"
	"testing"

	myerrors "github.com/ImmortaL-jsdev/task-manager/internal/errors"
	"github.com/ImmortaL-jsdev/task-manager/internal/models"
	"github.com/ImmortaL-jsdev/task-manager/internal/service"
	"github.com/ImmortaL-jsdev/task-manager/internal/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateTask_OK(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(models.Task{ID: "t1", Title: "write tests", Status: "todo"}, nil)

	svc := service.NewTaskService(repo)

	got, err := svc.CreateTask(context.Background(), models.Task{Title: "write tests", Status: "todo"})

	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if got.ID != "t1" {
		t.Fatalf("want ID t1, got %q", got.ID)
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	svc := service.NewTaskService(repo)

	_, err := svc.CreateTask(context.Background(), models.Task{Title: "", Status: "todo"})
	if err == nil {
		t.Fatalf("want error, got nil")
	}
	var valErr *myerrors.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("want *ValidationError, got %T", err)
	}
	if valErr.Message != "title cannot be empty" {
		t.Fatalf("want message %q, got %q", "title cannot be empty", valErr.Message)
	}
}

func TestCreateTask_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	svc := service.NewTaskService(repo)

	_, err := svc.CreateTask(context.Background(), models.Task{Title: "ok", Status: "archived"})
	if err == nil {
		t.Fatalf("want error, got nil")
	}

	var valErr *myerrors.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("want *ValidationError, got %T", err)
	}
	if valErr.Message != "invalid status" {
		t.Fatalf("want message %q, got %q", "invalid status", valErr.Message)
	}

}
func TestCreateTask_RepoError(t *testing.T) {
	sentinel := errors.New("db down")

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(models.Task{}, sentinel)

	svc := service.NewTaskService(repo)

	_, err := svc.CreateTask(context.Background(), models.Task{Title: "ok", Status: "todo"})

	if err == nil {
		t.Fatalf("want error, got nil")
	}

	if !errors.Is(err, sentinel) {
		t.Fatalf("want err to wrap sentinel , got %v", err)
	}
}
