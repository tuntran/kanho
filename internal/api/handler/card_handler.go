package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	mw "github.com/tungtran/kanho/internal/api/middleware"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/service"
)

type CardHandler struct {
	svc *service.CardService
}

func NewCardHandler(svc *service.CardService) *CardHandler {
	return &CardHandler{svc: svc}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid board ID")
		return
	}

	userID, ok := mw.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		ColumnID    string          `json:"column_id"`
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Priority    domain.Priority `json:"priority"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	colID, err := uuid.Parse(body.ColumnID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid column_id")
		return
	}

	card, err := h.svc.Create(r.Context(), service.CreateCardInput{
		BoardID:     boardID,
		ColumnID:    colID,
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
		CreatedBy:   userID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, card)
}

func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid board ID")
		return
	}

	cards, err := h.svc.ListByBoard(r.Context(), boardID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cards)
}

func (h *CardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	card, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "card not found")
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	userID, ok := mw.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		Title       *string          `json:"title"`
		Description *string          `json:"description"`
		Priority    *domain.Priority `json:"priority"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	card, err := h.svc.Update(r.Context(), id, service.UpdateCardInput{
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
	}, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (h *CardHandler) MoveCard(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	userID, ok := mw.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		ColumnID    string  `json:"column_id"`
		AfterCardID *string `json:"after_card_id"`
		BeforeCardID *string `json:"before_card_id"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	colID, err := uuid.Parse(body.ColumnID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid column_id")
		return
	}

	input := service.MoveCardInput{ColumnID: colID}
	if body.AfterCardID != nil {
		parsed, err := uuid.Parse(*body.AfterCardID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid after_card_id")
			return
		}
		input.AfterCardID = &parsed
	}
	if body.BeforeCardID != nil {
		parsed, err := uuid.Parse(*body.BeforeCardID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid before_card_id")
			return
		}
		input.BeforeCardID = &parsed
	}

	if err := h.svc.Move(r.Context(), id, input, userID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CardHandler) ArchiveCard(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	userID, ok := mw.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.svc.Archive(r.Context(), id, userID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	member := memberFromContext(r.Context())
	if member == nil || member.Role == domain.RoleMember {
		writeError(w, http.StatusForbidden, "only admin or owner can delete cards")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CardHandler) ListActivities(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	limit := 50
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	activities, err := h.svc.ListActivities(r.Context(), id, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, activities)
}

func (h *CardHandler) UpdateAssignees(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	var body struct {
		Add    []string `json:"add"`
		Remove []string `json:"remove"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	for _, uid := range body.Add {
		parsed, err := uuid.Parse(uid)
		if err != nil {
			continue
		}
		_ = h.svc.AddAssignee(r.Context(), id, parsed)
	}
	for _, uid := range body.Remove {
		parsed, err := uuid.Parse(uid)
		if err != nil {
			continue
		}
		_ = h.svc.RemoveAssignee(r.Context(), id, parsed)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CardHandler) UpdateLabels(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card ID")
		return
	}

	var body struct {
		Add    []string `json:"add"`
		Remove []string `json:"remove"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	for _, lid := range body.Add {
		parsed, err := uuid.Parse(lid)
		if err != nil {
			continue
		}
		_ = h.svc.AddLabel(r.Context(), id, parsed)
	}
	for _, lid := range body.Remove {
		parsed, err := uuid.Parse(lid)
		if err != nil {
			continue
		}
		_ = h.svc.RemoveLabel(r.Context(), id, parsed)
	}
	w.WriteHeader(http.StatusNoContent)
}
