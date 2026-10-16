package repository_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	myerrors "github.com/ImmortaL-jsdev/task-manager/internal/errors"
	"github.com/ImmortaL-jsdev/task-manager/internal/models"
	"github.com/ImmortaL-jsdev/task-manager/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testStore *repository.PostgresStore
	testPool  *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		panic(err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	testPool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	applyMigrations(ctx, testPool)

	testStore, err = repository.NewPostgresStore(dsn)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	testStore.Close()
	testPool.Close()
	_ = pgContainer.Terminate(ctx)

	os.Exit(code)

}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) {
	content, err := os.ReadFile("../../migrations/000001_create_tasks_table.up.sql")
	if err != nil {
		panic(err)
	}
	for _, stmt := range strings.Split(string(content), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := pool.Exec(ctx, stmt); err != nil {
			panic("migration failed: " + err.Error() + "\nStatement: " + stmt)
		}
	}
}

func cleanTasks(t *testing.T) {
	t.Helper()
	if _, err := testPool.Exec(context.Background(), "TRUNCATE tasks"); err != nil {
		t.Fatalf("failed to truncate tasks: %v", err)
	}
}

func TestPostgresStore_CreateTask(t *testing.T) {
	cleanTasks(t)

	ctx := context.Background()

	in := models.Task{
		Title:       "write tests",
		Description: "with testcontainer",
		Status:      "todo",
	}

	got, err := testStore.CreateTask(ctx, in)
	if err != nil {
		t.Fatalf("unexpexted error: %v", err)
	}

	if got.ID == "" {
		t.Fatalf("want non empty ID, got empty")
	}
	if got.Title != in.Title {
		t.Fatalf("want title %q got %q", in.Title, got.Title)
	}
	if got.Description != in.Description {
		t.Fatalf("want %q description got %q", in.Description, got.Title)
	}
	if got.Status != in.Status {
		t.Fatalf("want %q status, got %q", in.Status, got.Status)
	}
	if got.CreatedAt.IsZero() {
		t.Fatalf("want non zero CreatedAt")
	}
	if got.UpdatedAt.IsZero() {
		t.Fatalf("want non zero UpdatedAt")
	}
}

func TestPostgresStore_GetAllTasks(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	_, err := testStore.CreateTask(ctx, models.Task{
		Title:       "first task",
		Description: "description 1",
		Status:      "todo",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	_, err = testStore.CreateTask(ctx, models.Task{
		Title:       "second task",
		Description: "second description",
		Status:      "in_progress",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	got, err := testStore.GetAllTasks(ctx)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("want 2 tasks, got %d", len(got))
	}

	if got[0].Title != "second task" {
		t.Fatalf("want first title %q, got %q", "second task", got[0].Title)
	}
	if got[1].Title != "first task" {
		t.Fatalf("want second title %q, got %q", "first task", got[1].Title)

	}
	if got[0].ID == "" {
		t.Fatalf("want non empty ID for first task")
	}
	if got[0].CreatedAt.IsZero() {
		t.Fatalf("want non zero CreatedAt for first task")
	}
}
func TestPostgresStore_GetTaskByStatus(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	_, err := testStore.CreateTask(ctx, models.Task{
		Title:       "todo1",
		Description: "description 1",
		Status:      "todo",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	_, err = testStore.CreateTask(ctx, models.Task{
		Title:       "in_progress1",
		Description: "first in_progress description",
		Status:      "in_progress",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	_, err = testStore.CreateTask(ctx, models.Task{
		Title:       "todo2",
		Description: "second todo description",
		Status:      "todo",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	got, err := testStore.GetTaskByStatus(ctx, "todo")
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 tasks with status todo got %d", len(got))
	}

	if got[0].Title != "todo2" {
		t.Fatalf("want first title %q, got %q", "todo2", got[0].Title)
	}
	if got[1].Title != "todo1" {
		t.Fatalf("want second title %q, got %q", "todo1", got[1].Title)
	}

	if got[0].Status != "todo" {
		t.Fatalf("want status todo got %q", got[0].Status)
	}
	if got[1].Status != "todo" {
		t.Fatalf("want status todo, got %q", got[1].Status)
	}

	got, err = testStore.GetTaskByStatus(ctx, "done")
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}
	if len(got) != 0 {

		t.Fatalf("want 0 ,got %d", len(got))

	}
}

func TestPostgresStore_GetTaskByID_OK(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	created, err := testStore.CreateTask(ctx, models.Task{
		Title:       "first",
		Description: "firstdesc",
		Status:      "todo",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	got, err := testStore.GetTaskByID(ctx, created.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf("want id %q, got %q", created.ID, got.ID)
	}
	if got.Title != created.Title {
		t.Fatalf("want title %q, got %q", created.Title, got.Title)
	}
	if got.Status != created.Status {
		t.Fatalf("want status %q, got %q", created.Status, got.Status)
	}
	if got.CreatedAt.IsZero() {
		t.Fatalf("want non empty createdAt got empty")
	}
}
func TestPostgresStore_GetTaskByID_NotFound(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	_, err := testStore.GetTaskByID(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatalf("want error got nil")
	}

	var notFound *myerrors.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("want *NotFoundError, got %T", err)
	}

	if notFound.Entity != "task" {
		t.Fatalf("want entity %q, got %q", "task", notFound.Entity)
	}

	if notFound.ID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("want ID %q, got %q", "00000000-0000-0000-0000-000000000000", notFound.ID)
	}
}

func TestPostgresStore_Update_OK(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	created, err := testStore.CreateTask(ctx, models.Task{
		Title:       "first",
		Description: "firstdesc",
		Status:      "todo",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	time.Sleep(10 * time.Millisecond)

	got, err := testStore.UpdateTask(ctx, created.ID, models.Task{
		Title:       "updatedfirst",
		Description: "updatedfirstdesc",
		Status:      "done",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf("want %q id got %q", created.ID, got.ID)
	}
	if got.Title != "updatedfirst" {
		t.Fatalf("want %q title, got %q", "updatedfirst", got.Title)
	}
	if got.Description != "updatedfirstdesc" {
		t.Fatalf("want %q description, got %q", "updatedfirstdesc", got.Description)
	}
	if got.Status != "done" {
		t.Fatalf("want %q status, got %q", "done", got.Status)
	}
	if !got.CreatedAt.Equal(created.CreatedAt) {
		t.Fatalf("want no changed CreatedAt got changed")
	}
	if !got.UpdatedAt.After(created.UpdatedAt) {
		t.Fatalf("want changed updatedAt got no changed")
	}
}
func TestPostgresStore_Update_NotFound(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	_, err := testStore.UpdateTask(ctx, "00000000-0000-0000-0000-000000000000", models.Task{
		Title:  "any",
		Status: "todo",
	})
	if err == nil {
		t.Fatalf("want error got nil")
	}
	var notFound *myerrors.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("want *NotFoundError, got %T", err)
	}

	if notFound.Entity != "task" {
		t.Fatalf("want entity %q, got %q", "task", notFound.Entity)
	}

	if notFound.ID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("want ID %q, got %q", "00000000-0000-0000-0000-000000000000", notFound.ID)
	}
}

func TestPostgresStore_Delete_OK(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	created, err := testStore.CreateTask(ctx, models.Task{
		Title:       "first",
		Description: "firstdesc",
		Status:      "todo",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	err = testStore.DeleteTask(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = testStore.GetTaskByID(ctx, created.ID)
	if err == nil {
		t.Fatalf("want error after delete got nil")
	}

	var notFound *myerrors.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("want *NotFoundError after delete got %T", err)
	}
}
func TestPostgresStore_Delete_NotFound(t *testing.T) {
	cleanTasks(t)
	ctx := context.Background()

	err := testStore.DeleteTask(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatalf("want error got nil")
	}
	var notFound *myerrors.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("want *NotFoundError, got %T", err)
	}

	if notFound.Entity != "task" {
		t.Fatalf("want entity %q, got %q", "task", notFound.Entity)
	}

	if notFound.ID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("want ID %q, got %q", "00000000-0000-0000-0000-000000000000", notFound.ID)
	}
}
