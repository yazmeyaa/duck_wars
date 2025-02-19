package logger

import (
	"log/slog"
	"os"
)

func New(app string) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.Group("application_info",
			"app_name", app,
		),
	)

	return logger
}
