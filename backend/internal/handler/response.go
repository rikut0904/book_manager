package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func methodNotAllowed(w http.ResponseWriter, methods ...string) {
	if len(methods) > 0 {
		w.Header().Set("Allow", strings.Join(methods, ", "))
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
		"error": "method not allowed",
	})
}

func notFound(w http.ResponseWriter) {
	writeJSON(w, http.StatusNotFound, map[string]string{
		"error": "not found",
	})
}

func notFoundWithMessage(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusNotFound, map[string]string{
		"error":   "not found",
		"message": message,
	})
}

func badRequest(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{
		"error":   "bad request",
		"message": message,
	})
}

func unauthorized(w http.ResponseWriter) {
	writeJSON(w, http.StatusUnauthorized, map[string]string{
		"error": "unauthorized",
	})
}

func conflict(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusConflict, map[string]string{
		"error":   "conflict",
		"message": message,
	})
}

func internalError(w http.ResponseWriter) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error": "internal server error",
	})
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func echoRequest(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	payload := map[string]any{
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  r.URL.Query(),
		"body":   string(body),
	}
	var jsonBody any
	if len(body) > 0 && json.Unmarshal(body, &jsonBody) == nil {
		payload["json"] = jsonBody
	}
	writeJSON(w, http.StatusOK, payload)
}
