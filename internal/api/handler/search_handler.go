package handler

import (
	"net/http"

	"github.com/tungtran/kanho/internal/repository"
)

type SearchHandler struct {
	cardRepo *repository.CardRepo
}

func NewSearchHandler(cardRepo *repository.CardRepo) *SearchHandler {
	return &SearchHandler{cardRepo: cardRepo}
}

func (h *SearchHandler) SearchCards(w http.ResponseWriter, r *http.Request) {
	ws := workspaceFromContext(r.Context())
	if ws == nil {
		writeError(w, http.StatusNotFound, "workspace not found")
		return
	}

	query := r.URL.Query().Get("q")
	if len(query) < 2 {
		writeError(w, http.StatusBadRequest, "query must be at least 2 characters")
		return
	}

	results, err := h.cardRepo.SearchFullText(r.Context(), ws.ID, query, 20)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, results)
}
