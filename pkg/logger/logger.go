package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	return logger
}

func LogError(logger *slog.Logger, funcName string, msg string, err error) {
	logger.Error(fmt.Sprintf("[%s]: %s", funcName, msg), "error", err.Error())
}
