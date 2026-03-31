package handler

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// WriteErrorPublic is exported for use by middleware packages.
func WriteErrorPublic(w http.ResponseWriter, status int, msg string) {
	writeError(w, status, msg)
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}
