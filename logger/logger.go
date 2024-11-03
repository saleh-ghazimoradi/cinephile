// /internal/logger/logger.go
package logger

import (
	"log/slog"
	"os"
)

//var Logger *slog.Logger

var Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
