package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/cfegela/azure-aca-go-templ-mongo/internal/auth"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/database"
	"github.com/cfegela/azure-aca-go-templ-mongo/internal/models"
)

type AuthHandler struct {
	userRepo   *database.UserRepository
	inviteRepo *database.InviteRepository
	authConfig *auth.Config
	jwtExpiry  time.Duration
}

func NewAuthHandler(userRepo *database.UserRepository, inviteRepo *database.InviteRepository, authConfig *auth.Config, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
		authConfig: authConfig,
		jwtExpiry:  jwtExpiry,
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Redirect(w, r, "/login?error=missing_fields", http.StatusSeeOther)
		return
	}

	user, err := h.userRepo.FindByEmail(r.Context(), email)
	if err != nil {
		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.authConfig.JWTSecret, h.jwtExpiry)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.jwtExpiry.Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := strings.TrimPrefix(r.URL.Path, "/register/")
	if token == "" {
		http.Error(w, "Invalid invite token", http.StatusBadRequest)
		return
	}

	invite, err := h.inviteRepo.FindByToken(r.Context(), token)
	if err != nil {
		http.Redirect(w, r, "/register/"+token+"?error=invalid_invite", http.StatusSeeOther)
		return
	}

	if !invite.IsValid() {
		http.Redirect(w, r, "/register/"+token+"?error=invite_expired", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if name == "" || email == "" || password == "" {
		http.Redirect(w, r, "/register/"+token+"?error=missing_fields", http.StatusSeeOther)
		return
	}

	if email != invite.Email {
		http.Redirect(w, r, "/register/"+token+"?error=email_mismatch", http.StatusSeeOther)
		return
	}

	if password != confirmPassword {
		http.Redirect(w, r, "/register/"+token+"?error=password_mismatch", http.StatusSeeOther)
		return
	}

	if err := models.ValidatePassword(password); err != nil {
		http.Redirect(w, r, "/register/"+token+"?error="+err.Error(), http.StatusSeeOther)
		return
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         name,
		Role:         models.RoleUser,
	}

	if err := user.Validate(); err != nil {
		http.Redirect(w, r, "/register/"+token+"?error="+err.Error(), http.StatusSeeOther)
		return
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		http.Redirect(w, r, "/register/"+token+"?error="+err.Error(), http.StatusSeeOther)
		return
	}

	if err := h.inviteRepo.MarkUsed(r.Context(), token); err != nil {
		// User created but invite not marked - log this but continue
	}

	jwtToken, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.authConfig.JWTSecret, h.jwtExpiry)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.jwtExpiry.Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
