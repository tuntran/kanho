package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/service"
)

type ColumnHandler struct {
	svc *service.ColumnService
}

func NewColumnHandler(svc *service.ColumnService) *ColumnHandler {
	return &ColumnHandler{svc: svc}
}

func (h *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Only admin/owner can create columns
	member := memberFromContext(r.Context())
	if member == nil || member.Role == domain.RoleMember {
		writeError(w, http.StatusForbidden, "only admin or owner can create columns")
		return
	}

	boardID, err := uuid.Parse(chi.URLParam(r, "boardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid board ID")
		return
	}

	var body struct {
		Name     string `json:"name"`
		Color    string `json:"color"`
		WIPLimit *int   `json:"wip_limit"`
		IsDone   bool   `json:"is_done"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	col, err := h.svc.Create(r.Context(), service.CreateColumnInput{
		BoardID:  boardID,
		Name:     body.Name,
		Color:    body.Color,
		WIPLimit: body.WIPLimit,
		IsDone:   body.IsDone,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, col)
}

func (h *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Only admin/owner can edit columns
	member := memberFromContext(r.Context())
	if member == nil || member.Role == domain.RoleMember {
		writeError(w, http.StatusForbidden, "only admin or owner can edit columns")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "columnID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid column ID")
		return
	}

	var body struct {
		Name     string `json:"name"`
		Color    string `json:"color"`
		WIPLimit *int   `json:"wip_limit"`
		IsDone   *bool  `json:"is_done"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	col, err := h.svc.Update(r.Context(), id, service.UpdateColumnInput{
		Name:     body.Name,
		Color:    body.Color,
		WIPLimit: body.WIPLimit,
		IsDone:   body.IsDone,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, col)
}

func (h *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Only admin/owner can delete columns
	member := memberFromContext(r.Context())
	if member == nil || member.Role == domain.RoleMember {
		writeError(w, http.StatusForbidden, "only admin or owner can delete columns")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "columnID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid column ID")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ColumnHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	// Only admin/owner can reorder columns
	member := memberFromContext(r.Context())
	if member == nil || member.Role == domain.RoleMember {
		writeError(w, http.StatusForbidden, "only admin or owner can reorder columns")
		return
	}

	var body []service.ReorderItem
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Reorder(r.Context(), body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
