package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	myerrors "github.com/ImmortaL-jsdev/task-manager/internal/errors"
	"github.com/ImmortaL-jsdev/task-manager/internal/models"
	"github.com/ImmortaL-jsdev/task-manager/internal/service"
	"github.com/gorilla/mux"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tasks, err := h.service.GetAllTasks(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	created, err := h.service.CreateTask(ctx, task)
	if err != nil {
		var valError *myerrors.ValidationError

		if errors.As(err, &valError) {
			respondWithError(w, http.StatusBadRequest, valError.Message)
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	task, err := h.service.GetTaskByID(ctx, id)
	if err != nil {
		var notFound *myerrors.NotFoundError
		if errors.As(err, &notFound) {
			respondWithError(w, http.StatusNotFound, notFound.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	vars := mux.Vars(r)
	id := vars["id"]

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	updated, err := h.service.UpdateTask(ctx, id, task)

	if err != nil {
		var notFound *myerrors.NotFoundError
		if errors.As(err, &notFound) {
			respondWithError(w, http.StatusNotFound, notFound.Error())
			return
		}
		var valErr *myerrors.ValidationError
		if errors.As(err, &valErr) {
			respondWithError(w, http.StatusBadRequest, valErr.Message)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respondWithJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.service.DeleteTask(ctx, id)

	if err != nil {
		var notFound *myerrors.NotFoundError
		if errors.As(err, &notFound) {
			respondWithError(w, http.StatusNotFound, notFound.Error())
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
