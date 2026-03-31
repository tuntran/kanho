package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	mw "github.com/tungtran/kanho/internal/api/middleware"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/service"
)

type WorkspaceHandler struct {
	svc *service.WorkspaceService
}

func NewWorkspaceHandler(svc *service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{svc: svc}
}

func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := mw.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	workspaces, err := h.svc.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list workspaces")
		return
	}
	writeJSON(w, http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := mw.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	ws, err := h.svc.Create(r.Context(), service.CreateWorkspaceInput{
		Name:        body.Name,
		Description: body.Description,
		CreatedBy:   userID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, ws)
}

func (h *WorkspaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "workspaceSlug")
	ws, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}
	writeJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		LogoURL     string `json:"logo_url"`
		AccentColor string `json:"accent_color"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.svc.Update(r.Context(), ws.ID, service.UpdateWorkspaceInput{
		Name:        body.Name,
		Description: body.Description,
		LogoURL:     body.LogoURL,
		AccentColor: body.AccentColor,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *WorkspaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	member := memberFromContext(r.Context())
	if member == nil || member.Role != domain.RoleOwner {
		writeError(w, http.StatusForbidden, "only owner can delete workspace")
		return
	}

	if err := h.svc.Delete(r.Context(), ws.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	members, err := h.svc.ListMembers(r.Context(), ws.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, members)
}

func (h *WorkspaceHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	member := memberFromContext(r.Context())
	if member == nil || member.Role == domain.RoleMember {
		writeError(w, http.StatusForbidden, "only admin or owner can invite members")
		return
	}

	var body struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	uid, err := uuid.Parse(body.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	role := domain.Role(body.Role)
	if role == "" {
		role = domain.RoleMember
	}

	if err := h.svc.InviteMember(r.Context(), service.InviteMemberInput{
		WorkspaceID: ws.ID,
		UserID:      uid,
		Role:        role,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Context helpers — read values set by workspace middleware

func workspaceFromContext(ctx interface{ Value(any) any }) *domain.Workspace {
	if v, ok := ctx.Value(mw.WorkspaceKey).(*domain.Workspace); ok {
		return v
	}
	return nil
}

func memberFromContext(ctx interface{ Value(any) any }) *domain.WorkspaceMember {
	if v, ok := ctx.Value(mw.MemberKey).(*domain.WorkspaceMember); ok {
		return v
	}
	return nil
}
