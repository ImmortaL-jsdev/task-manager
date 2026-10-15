package service_test

import (
	"context"
	"testing"

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
	t.Logf("got error: %v (type %T)", err, err)
}
