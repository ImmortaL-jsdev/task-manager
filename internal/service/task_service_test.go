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

func TestGetTaskByStatus_EmptyStatus_CallsGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().GetAllTasks(gomock.Any()).Return([]models.Task{{ID: "t1"}, {ID: "t2"}}, nil)

	svc := service.NewTaskService(repo)

	got, err := svc.GetTaskByStatus(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 tasks, got %d", len(got))
	}

}
func TestGetTaskByStatus_ValidStatus_CallsGetByStatus(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)
	repo.EXPECT().GetTaskByStatus(gomock.Any(), "todo").Return([]models.Task{{ID: "1", Status: "todo"}}, nil)

	svc := service.NewTaskService(repo)

	got, err := svc.GetTaskByStatus(context.Background(), "todo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 tasks, got %d", len(got))
	}
	if got[0].Status != "todo" {
		t.Fatalf("want status todo, got %q", got[0].Status)
	}
}

func TestGetTaskByStatus_InvalidStatus_RepoNotCalled(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)
	svc := service.NewTaskService(repo)

	got, err := svc.GetTaskByStatus(context.Background(), "archived")

	if err == nil {
		t.Fatalf("want error, got nil")
	}

	if got != nil {
		t.Fatalf("want nil slice, got %v", got)
	}

	var valErr *myerrors.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("want ValidationError, got %T", err)
	}
	if valErr.Message != "invalid status" {
		t.Fatalf("want message %q, got %q", "invalid status", valErr.Message)
	}
}

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().GetAllTasks(gomock.Any()).Return([]models.Task{{ID: "t1"}, {ID: "t2"}}, nil)

	svc := service.NewTaskService(repo)

	got, err := svc.GetAllTasks(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("want 2 tasks, got %d", len(got))
	}

	if got[0].ID != "t1" {
		t.Fatalf("want first ID t1, got %q", got[0].ID)
	}
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().GetTaskByID(gomock.Any(), "t1").Return(models.Task{ID: "t1"}, nil)

	svc := service.NewTaskService(repo)

	got, err := svc.GetTaskByID(context.Background(), "t1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != "t1" {
		t.Fatalf("want ID t1, got %q", got.ID)
	}
}

func TestDeleteTaskOK(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().DeleteTask(gomock.Any(), "t1").Return(nil)

	svc := service.NewTaskService(repo)

	err := svc.DeleteTask(context.Background(), "t1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestDeleteTask_RepoError(t *testing.T) {
	sentinel := errors.New("not found")

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)
	repo.EXPECT().DeleteTask(gomock.Any(), "t1").Return(sentinel)

	svc := service.NewTaskService(repo)

	err := svc.DeleteTask(context.Background(), "t1")
	if err == nil {
		t.Fatalf("want error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("want err to wrap sentinel, got %v", err)
	}

}
func TestUpdateTask_OK(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)
	repo.EXPECT().UpdateTask(gomock.Any(), "t1", gomock.Any()).Return(models.Task{ID: "t1", Title: "updated", Status: "in_progress"}, nil)

	svc := service.NewTaskService(repo)

	got, err := svc.UpdateTask(context.Background(), "t1", models.Task{Title: "updated", Status: "in_progress"})

	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if got.ID != "t1" {
		t.Fatalf("want ID t1, got %q", got.ID)
	}
	if got.Title != "updated" {
		t.Fatalf("want title updated, got %q", got.Title)
	}
}

func TestUpdateTask_DefaultStatus(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)
	repo.EXPECT().UpdateTask(gomock.Any(), "t1", gomock.Any()).DoAndReturn(func(_ context.Context, id string, task models.Task) (models.Task, error) {
		if task.Status != "todo" {
			t.Fatalf("repo got status %q, want todo", task.Status)
		}
		task.ID = id

		return task, nil
	})
	svc := service.NewTaskService(repo)

	got, err := svc.UpdateTask(context.Background(), "t1", models.Task{Title: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != "todo" {
		t.Fatalf("want status todo, got %q", got.Status)
	}
}

func TestUpdateTask_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	svc := service.NewTaskService(repo)

	_, err := svc.UpdateTask(context.Background(), "t1", models.Task{Title: "", Status: "todo"})
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

func TestUpdateTask_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	svc := service.NewTaskService(repo)

	_, err := svc.UpdateTask(context.Background(), "t1", models.Task{Title: "ok", Status: "archived"})
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

func TestUpdateTask_RepoError(t *testing.T) {
	sentinel := errors.New("update failed")

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockTaskRepository(ctrl)

	repo.EXPECT().UpdateTask(gomock.Any(), "t1", gomock.Any()).Return(models.Task{}, sentinel)

	svc := service.NewTaskService(repo)

	_, err := svc.UpdateTask(context.Background(), "t1", models.Task{Title: "ok", Status: "todo"})

	if err == nil {
		t.Fatalf("want error, got nil")
	}

	if !errors.Is(err, sentinel) {
		t.Fatalf("failed to update task: %v", err)
	}
}
