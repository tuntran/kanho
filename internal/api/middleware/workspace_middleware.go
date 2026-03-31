package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tungtran/kanho/internal/service"
)

type workspaceCtxKey string

const (
	WorkspaceKey workspaceCtxKey = "workspace"
	MemberKey    workspaceCtxKey = "member"
)

// WorkspaceMiddleware loads workspace from slug and checks membership.
func WorkspaceMiddleware(workspaceSvc *service.WorkspaceService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slug := chi.URLParam(r, "workspaceSlug")
			if slug == "" {
				writeMiddlewareError(w, http.StatusBadRequest, "workspace slug required")
				return
			}

			ws, err := workspaceSvc.GetBySlug(r.Context(), slug)
			if err != nil {
				writeMiddlewareError(w, http.StatusNotFound, "workspace not found")
				return
			}

			ctx := context.WithValue(r.Context(), WorkspaceKey, ws)

			// Check membership
			userID, ok := GetUserID(r.Context())
			if !ok {
				writeMiddlewareError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			member, err := workspaceSvc.GetMember(r.Context(), ws.ID, userID)
			if err != nil {
				writeMiddlewareError(w, http.StatusForbidden, "not a member of this workspace")
				return
			}
			ctx = context.WithValue(ctx, MemberKey, member)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeMiddlewareError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
