package transport_http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

// Write an HTTP error wtith provided status code and message in JSON
// {"error": errorMsg}
func WriteError(w http.ResponseWriter, code int, errorMsg string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(map[string]string{
		"error": errorMsg,
	})
}

func handleParseError(w http.ResponseWriter, logger *slog.Logger, fn string, err error, message string, statusCode int) error {
	logger.Error(message, "fn", fn, "error", err.Error())
	WriteError(w, statusCode, message)
	return err
}

// Parse any JSON struct. Provide only reference to structs as dest
func ParseJSON(dest any, w http.ResponseWriter, r *http.Request, logger *slog.Logger, fn string) error {
	err := json.NewDecoder(r.Body).Decode(&dest)
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error("Cannot close request body.")
		}
	}()

	if errors.Is(err, io.EOF) {
		return handleParseError(w, logger, fn, err, "Received empty request body", http.StatusBadRequest)
	} else if err != nil {
		return handleParseError(w, logger, fn, err, "Cannot parse request body.", http.StatusBadRequest)
	}

	return nil
}
