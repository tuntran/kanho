package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tungtran/kanho/internal/storage"
)

const (
	maxFileSize       = 25 << 20 // 25 MB
	presignedPutTTL   = 15 * time.Minute
	uploadsBucketPath = "attachments"
)

var allowedContentTypes = map[string]bool{
	"image/png":      true,
	"image/jpeg":     true,
	"image/gif":      true,
	"image/webp":     true,
	"image/svg+xml":  true,
	"application/pdf": true,
}

type UploadHandler struct {
	store     storage.Client
	publicURL string
}

func NewUploadHandler(store storage.Client, publicURL string) *UploadHandler {
	return &UploadHandler{store: store, publicURL: publicURL}
}

func (h *UploadHandler) PresignUpload(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		CardID      string `json:"card_id"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Filename == "" || body.ContentType == "" {
		writeError(w, http.StatusBadRequest, "filename and content_type are required")
		return
	}

	if !allowedContentTypes[body.ContentType] {
		writeError(w, http.StatusBadRequest, "unsupported content type")
		return
	}

	// Generate UUID-based key to prevent path traversal
	ext := filepath.Ext(body.Filename)
	if ext == "" {
		ext = mimeToExt(body.ContentType)
	}
	fileKey := fmt.Sprintf("%s/%s%s", uploadsBucketPath, uuid.New().String(), ext)

	presignedURL, err := h.store.GetPresignedPutURL(r.Context(), fileKey, presignedPutTTL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate upload URL")
		return
	}

	publicFileURL := h.publicURL + "/" + fileKey

	writeJSON(w, http.StatusOK, map[string]string{
		"upload_url": presignedURL,
		"file_key":   fileKey,
		"public_url": publicFileURL,
	})
}

func mimeToExt(contentType string) string {
	switch strings.ToLower(contentType) {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	case "application/pdf":
		return ".pdf"
	default:
		return ""
	}
}
