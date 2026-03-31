package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/service"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	projects, err := h.svc.ListByWorkspace(r.Context(), ws.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	// Only admin/owner can create projects
	member := memberFromContext(r.Context())
	if member == nil || member.Role == domain.RoleMember {
		writeError(w, http.StatusForbidden, "only admin or owner can create projects")
		return
	}

	var body struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" || body.Key == "" {
		writeError(w, http.StatusBadRequest, "name and key are required")
		return
	}

	p, err := h.svc.Create(r.Context(), service.CreateProjectInput{
		WorkspaceID: ws.ID,
		Name:        body.Name,
		Key:         body.Key,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	var body struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.svc.Update(r.Context(), id, service.UpdateProjectInput{
		Name: body.Name,
		Key:  body.Key,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Only owner can delete
	member := memberFromContext(r.Context())
	if member == nil || member.Role != domain.RoleOwner {
		writeError(w, http.StatusForbidden, "only owner can delete projects")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
