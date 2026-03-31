package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/service"
)

type BoardHandler struct {
	boardSvc  *service.BoardService
	columnSvc *service.ColumnService
}

func NewBoardHandler(boardSvc *service.BoardService, columnSvc *service.ColumnService) *BoardHandler {
	return &BoardHandler{boardSvc: boardSvc, columnSvc: columnSvc}
}

func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "boardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid board ID")
		return
	}

	board, err := h.boardSvc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "board not found")
		return
	}

	columns, err := h.columnSvc.ListByBoard(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load columns")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"board":   board,
		"columns": columns,
	})
}

func (h *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "boardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid board ID")
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

	board, err := h.boardSvc.Update(r.Context(), id, body.Name, body.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, board)
}

func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "boardID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid board ID")
		return
	}

	if err := h.boardSvc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BoardHandler) ListBoards(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	boards, err := h.boardSvc.ListByProject(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, boards)
}
