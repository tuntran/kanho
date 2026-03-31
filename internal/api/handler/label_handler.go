package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/service"
)

type LabelHandler struct {
	svc *service.LabelService
}

func NewLabelHandler(svc *service.LabelService) *LabelHandler {
	return &LabelHandler{svc: svc}
}

func (h *LabelHandler) List(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	labels, err := h.svc.ListByWorkspace(r.Context(), ws.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, labels)
}

func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	var body struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	label, err := h.svc.Create(r.Context(), service.CreateLabelInput{
		WorkspaceID: ws.ID,
		Name:        body.Name,
		Color:       body.Color,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, label)
}

func (h *LabelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "labelID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid label ID")
		return
	}

	var body struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	label, err := h.svc.Update(r.Context(), id, service.UpdateLabelInput{
		Name:  body.Name,
		Color: body.Color,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, label)
}

func (h *LabelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "labelID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid label ID")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
