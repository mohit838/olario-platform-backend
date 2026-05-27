package httpapi

import (
	"encoding/json"
	"net/http"

	authapp "github.com/mohit838/olario-platform-backend/internal/application/auth"
)

type AuthHandler struct {
	service *authapp.Service
}

type loginRequest struct {
	TenantSlug string `json:"tenant_slug"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type registerRequest struct {
	InvitationToken string `json:"invitation_token"`
	TenantName      string `json:"tenant_name"`
	TenantSlug      string `json:"tenant_slug"`
	AdminName       string `json:"admin_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewAuthHandler(service *authapp.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth service unavailable"})
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	tokens, err := h.service.Register(r.Context(), authapp.RegisterInput{
		InvitationToken: req.InvitationToken,
		TenantName:      req.TenantName,
		TenantSlug:      req.TenantSlug,
		AdminName:       req.AdminName,
		Email:           req.Email,
		Password:        req.Password,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth service unavailable"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	if req.TenantSlug == "" || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tenant_slug, email, and password are required"})
		return
	}

	tokens, err := h.service.Login(r.Context(), authapp.LoginInput{
		TenantSlug: req.TenantSlug,
		Email:      req.Email,
		Password:   req.Password,
	})
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid login credentials"})
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth service unavailable"})
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh_token is required"})
		return
	}

	tokens, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth service unavailable"})
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "refresh_token is required"})
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "logout failed"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
