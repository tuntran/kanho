package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tungtran/kanho/internal/api"
	"github.com/tungtran/kanho/internal/api/ws"
	"github.com/tungtran/kanho/internal/config"
	"github.com/tungtran/kanho/internal/repository"
	"github.com/tungtran/kanho/internal/storage"

	kanho "github.com/tungtran/kanho"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Run DB migrations
	if cfg.Database.URL != "" {
		if err := repository.RunMigrations(cfg.Database.URL); err != nil {
			slog.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Init storage client
	storageClient, err := initStorage(cfg)
	if err != nil {
		slog.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	// Embedded frontend
	staticFS, err := fs.Sub(kanho.StaticFiles, "web/dist")
	if err != nil {
		slog.Error("failed to load static files", "error", err)
		os.Exit(1)
	}

	// WebSocket hub
	wsHub := ws.NewHub(pool)
	wsHub.Start(ctx)

	router := api.NewRouter(pool, storageClient, staticFS, cfg, wsHub)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
}

func initStorage(cfg *config.Config) (storage.Client, error) {
	if cfg.Storage.Endpoint == "" {
		return storage.NewLocalClient("./uploads")
	}

	client, err := storage.NewMinIOClient(
		cfg.Storage.Endpoint,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.Bucket,
		cfg.Storage.UseSSL,
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.EnsureBucket(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
