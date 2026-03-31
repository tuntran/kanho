package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
	"github.com/tungtran/kanho/internal/service"
)

type boardCtxKey string

const BoardKey boardCtxKey = "board"

// roleLevel maps roles to numeric levels for comparison.
var roleLevel = map[domain.Role]int{
	domain.RoleMember: 1,
	domain.RoleAdmin:  2,
	domain.RoleOwner:  3,
}

// BoardMiddleware loads board from {boardID} URL param, injects board and workspace
// membership into context so RequireRole can function on direct-board routes.
func BoardMiddleware(boardSvc *service.BoardService, boardRepo *repository.BoardRepo, workspaceSvc *service.WorkspaceService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			boardIDStr := chi.URLParam(r, "boardID")
			if boardIDStr == "" {
				writeMiddlewareError(w, http.StatusBadRequest, "board ID required")
				return
			}

			boardID, err := uuid.Parse(boardIDStr)
			if err != nil {
				writeMiddlewareError(w, http.StatusBadRequest, "invalid board ID")
				return
			}

			board, err := boardSvc.GetByID(r.Context(), boardID)
			if err != nil {
				writeMiddlewareError(w, http.StatusNotFound, "board not found")
				return
			}

			ctx := context.WithValue(r.Context(), BoardKey, board)

			// Load workspace membership so RequireRole works on direct-board routes.
			userID, ok := GetUserID(r.Context())
			if ok {
				workspaceID, err := boardRepo.GetWorkspaceIDByBoardID(r.Context(), boardID)
				if err == nil {
					member, err := workspaceSvc.GetMember(r.Context(), workspaceID, userID)
					if err == nil {
						ctx = context.WithValue(ctx, MemberKey, member)
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole enforces that the workspace member (from context) has at least minRole.
// Must be used after WorkspaceMiddleware or BoardMiddleware which sets MemberKey in context.
func RequireRole(minRole domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			member, ok := r.Context().Value(MemberKey).(*domain.WorkspaceMember)
			if !ok || member == nil {
				writeMiddlewareError(w, http.StatusForbidden, "membership required")
				return
			}

			memberLvl := roleLevel[member.Role]
			requiredLvl := roleLevel[minRole]

			if memberLvl < requiredLvl {
				writeMiddlewareError(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetBoard retrieves the board from request context.
func GetBoard(ctx context.Context) (*domain.Board, bool) {
	b, ok := ctx.Value(BoardKey).(*domain.Board)
	return b, ok
}
