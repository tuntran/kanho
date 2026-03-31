package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/api/handler"
	mw "github.com/tungtran/kanho/internal/api/middleware"
	"github.com/tungtran/kanho/internal/api/ws"
	"github.com/tungtran/kanho/internal/config"
	"github.com/tungtran/kanho/internal/domain"
	"github.com/tungtran/kanho/internal/repository"
	"github.com/tungtran/kanho/internal/service"
	"github.com/tungtran/kanho/internal/storage"
)

func NewRouter(db *pgxpool.Pool, storageClient storage.Client, staticFS fs.FS, cfg *config.Config, wsHub *ws.Hub) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(mw.CORS())

	// Repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	// Handlers
	healthHandler := handler.NewHealthHandler(db, storageClient)
	authHandler := handler.NewAuthHandler(authSvc, cfg.JWT.RefreshTTL, false)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler.ServeHTTP)

		r.Route("/v1", func(r chi.Router) {
			// Public auth routes
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", authHandler.Register)
				r.Post("/login", authHandler.Login)
				r.Post("/logout", authHandler.Logout)
				r.Post("/refresh", authHandler.Refresh)

				// Protected auth routes
				r.Group(func(r chi.Router) {
					r.Use(mw.JWTAuth(cfg.JWT.Secret))
					r.Get("/me", authHandler.Me)
				})
			})

			r.Get("/setup/status", authHandler.SetupStatus)

			// Protected API routes
			r.Group(func(r chi.Router) {
				r.Use(mw.JWTAuth(cfg.JWT.Secret))

				// --- Repositories ---
				workspaceRepo := repository.NewWorkspaceRepo(db)
				projectRepo := repository.NewProjectRepo(db)
				boardRepo := repository.NewBoardRepo(db)
				columnRepo := repository.NewColumnRepo(db)
				labelRepo := repository.NewLabelRepo(db)
				cardRepo := repository.NewCardRepo(db)
				activityRepo := repository.NewActivityRepo(db)

				commentRepo := repository.NewCommentRepo(db)

				// --- Services ---
				workspaceSvc := service.NewWorkspaceService(workspaceRepo)
				boardSvc := service.NewBoardService(boardRepo, columnRepo)
				projectSvc := service.NewProjectService(projectRepo, boardSvc)
				columnSvc := service.NewColumnService(columnRepo)
				labelSvc := service.NewLabelService(labelRepo)
				cardSvc := service.NewCardService(cardRepo, activityRepo, projectRepo, boardRepo)
				commentSvc := service.NewCommentService(commentRepo)

				// --- Handlers ---
				workspaceHandler := handler.NewWorkspaceHandler(workspaceSvc)
				projectHandler := handler.NewProjectHandler(projectSvc)
				boardHandler := handler.NewBoardHandler(boardSvc, columnSvc)
				columnHandler := handler.NewColumnHandler(columnSvc)
				labelHandler := handler.NewLabelHandler(labelSvc)
				cardHandler := handler.NewCardHandler(cardSvc)
				commentHandler := handler.NewCommentHandler(commentSvc)
				searchHandler := handler.NewSearchHandler(cardRepo)
				wsHandler := ws.NewHandler(wsHub, cfg.JWT.Secret)

				// Workspaces
				r.Route("/workspaces", func(r chi.Router) {
					r.Get("/", workspaceHandler.List)
					r.Post("/", workspaceHandler.Create)
					r.Route("/{workspaceSlug}", func(r chi.Router) {
						r.Use(mw.WorkspaceMiddleware(workspaceSvc))
						r.Get("/", workspaceHandler.Get)
						r.Patch("/", workspaceHandler.Update)
						r.Delete("/", workspaceHandler.Delete)

						// Members
						r.Get("/members", workspaceHandler.ListMembers)
						r.Post("/members", workspaceHandler.InviteMember)

						// Projects
						r.Route("/projects", func(r chi.Router) {
							r.Get("/", projectHandler.List)
							r.Post("/", projectHandler.Create)
							r.Route("/{projectID}", func(r chi.Router) {
								r.Get("/", projectHandler.Get)
								r.Patch("/", projectHandler.Update)
								r.Delete("/", projectHandler.Delete)
								r.Get("/boards", boardHandler.ListBoards)
							})
						})

						// Labels
						r.Route("/labels", func(r chi.Router) {
							r.Get("/", labelHandler.List)
							r.Post("/", labelHandler.Create)
							r.Route("/{labelID}", func(r chi.Router) {
								r.Patch("/", labelHandler.Update)
								r.Delete("/", labelHandler.Delete)
							})
						})

						// Search
						r.Get("/search", searchHandler.SearchCards)
					})
				})

				// Boards (direct access by ID)
				r.Route("/boards/{boardID}", func(r chi.Router) {
					r.Use(mw.BoardMiddleware(boardSvc, boardRepo, workspaceSvc))
					r.Get("/", boardHandler.Get)
					r.With(mw.RequireRole(domain.RoleAdmin)).Patch("/", boardHandler.Update)
					r.With(mw.RequireRole(domain.RoleOwner)).Delete("/", boardHandler.Delete)

					// Column create requires admin role
					r.With(mw.RequireRole(domain.RoleAdmin)).Post("/columns", columnHandler.Create)

					// Cards
					r.Get("/cards", cardHandler.List)
					r.Post("/cards", cardHandler.Create)

					// WebSocket
					r.Get("/ws", wsHandler.ServeHTTP)
				})

				// Attachments
				uploadHandler := handler.NewUploadHandler(storageClient, cfg.Storage.PublicURL)
				r.Post("/attachments/upload-url", uploadHandler.PresignUpload)

				// Cards (direct access by ID)
				r.Route("/cards/{cardID}", func(r chi.Router) {
					r.Get("/", cardHandler.Get)
					r.Patch("/", cardHandler.Update)
					r.Delete("/", cardHandler.DeleteCard)
					r.Post("/move", cardHandler.MoveCard)
					r.Post("/archive", cardHandler.ArchiveCard)
					r.Patch("/assignees", cardHandler.UpdateAssignees)
					r.Patch("/labels", cardHandler.UpdateLabels)
					r.Get("/activity", cardHandler.ListActivities)

					// Comments
					r.Get("/comments", commentHandler.List)
					r.Post("/comments", commentHandler.Create)
				})

				// Comments (direct access by ID)
				r.Route("/comments/{commentID}", func(r chi.Router) {
					r.Patch("/", commentHandler.Update)
					r.Delete("/", commentHandler.Delete)
				})

				// Columns (direct access by ID — role enforced at board level for create)
				r.Route("/columns", func(r chi.Router) {
					r.Patch("/reorder", columnHandler.Reorder)
					r.Route("/{columnID}", func(r chi.Router) {
						r.Patch("/", columnHandler.Update)
						r.Delete("/", columnHandler.Delete)
					})
				})
			})
		})
	})

	// SPA fallback: serve static files, fallback to index.html
	spaHandler := spaFileServer(staticFS)
	r.Handle("/*", spaHandler)

	return r
}

func spaFileServer(staticFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(staticFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try to serve the file directly
		if _, err := fs.Stat(staticFS, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for client-side routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
