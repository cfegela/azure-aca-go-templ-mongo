package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cfegela/azure-aca-go-templ-mongo/internal/auth"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/database"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/models"
)

type TaskHandler struct {
	repo *database.TaskRepository
}

func NewTaskHandler(repo *database.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/api/tasks/") {
			h.GetTask(w, r)
		} else {
			h.ListTasks(w, r)
		}
	case http.MethodPost:
		h.CreateTask(w, r)
	case http.MethodPut:
		h.UpdateTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tasks, err := h.repo.FindByUserID(r.Context(), claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	task, err := h.repo.FindByIDAndUserID(r.Context(), id, claims.UserID)
	if err != nil {
		if err.Error() == "task not found" || err.Error() == "invalid task ID" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	task.UserID = claims.UserID

	if err := task.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Create(r.Context(), &task); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := task.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.UpdateByUserID(r.Context(), id, claims.UserID, &task); err != nil {
		if err.Error() == "task not found" || err.Error() == "invalid task ID" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	updatedTask, err := h.repo.FindByIDAndUserID(r.Context(), id, claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	if err := h.repo.DeleteByUserID(r.Context(), id, claims.UserID); err != nil {
		if err.Error() == "task not found" || err.Error() == "invalid task ID" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
