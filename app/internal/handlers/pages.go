package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/cfegela/azure-aca-go-templ-mongo/internal/auth"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/database"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/models"
	"github.com/cfegela/azure-aca-go-templ-mongo/web/templates"
)

type PageHandler struct {
	taskRepo   *database.TaskRepository
	userRepo   *database.UserRepository
	inviteRepo *database.InviteRepository
}

func NewPageHandler(taskRepo *database.TaskRepository, userRepo *database.UserRepository, inviteRepo *database.InviteRepository) *PageHandler {
	return &PageHandler{
		taskRepo:   taskRepo,
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
	}
}

func (h *PageHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")
	templates.Login(errorMsg).Render(r.Context(), w)
}

func (h *PageHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/register/")
	errorMsg := r.URL.Query().Get("error")

	invite, err := h.inviteRepo.FindByToken(r.Context(), token)
	if err != nil || !invite.IsValid() {
		http.Error(w, "Invalid or expired invite", http.StatusBadRequest)
		return
	}

	templates.Register(token, invite.Email, errorMsg).Render(r.Context(), w)
}

func (h *PageHandler) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tasks, err := h.taskRepo.FindByUserID(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

	templates.Dashboard(claims.Email, tasks).Render(r.Context(), w)
}

func (h *PageHandler) ShowTaskForm(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	templates.TaskForm(claims.Email, nil, false).Render(r.Context(), w)
}

func (h *PageHandler) ShowEditForm(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id = strings.TrimSuffix(id, "/edit")

	task, err := h.taskRepo.FindByIDAndUserID(r.Context(), id, claims.UserID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	templates.TaskForm(claims.Email, task, true).Render(r.Context(), w)
}

func (h *PageHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	task := &models.Task{
		UserID:      claims.UserID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
	}

	dueDateStr := r.FormValue("due_date")
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			task.DueDate = &dueDate
		}
	}

	if err := task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.taskRepo.Create(r.Context(), task); err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PageHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/tasks/")

	task := &models.Task{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
	}

	dueDateStr := r.FormValue("due_date")
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			task.DueDate = &dueDate
		}
	}

	if err := task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.taskRepo.UpdateByUserID(r.Context(), id, claims.UserID, task); err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PageHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id = strings.TrimSuffix(id, "/delete")

	if err := h.taskRepo.DeleteByUserID(r.Context(), id, claims.UserID); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PageHandler) ShowInvites(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	invites, err := h.inviteRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to load invites", http.StatusInternalServerError)
		return
	}

	successMsg := r.URL.Query().Get("success")
	templates.Invites(claims.Email, invites, successMsg).Render(r.Context(), w)
}

func (h *PageHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/admin/invites?error=missing_email", http.StatusSeeOther)
		return
	}

	token, err := models.GenerateInviteToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	invite := &models.Invite{
		Token:     token,
		Email:     email,
		InvitedBy: claims.UserID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := h.inviteRepo.Create(r.Context(), invite); err != nil {
		http.Error(w, "Failed to create invite", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/invites?success=invite_created", http.StatusSeeOther)
}
