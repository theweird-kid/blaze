package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/theweird-kid/blaze/internal/logger"
	"github.com/theweird-kid/blaze/internal/worker"
)

func main() {
	logger.Setup()

	mux := http.NewServeMux()
	mux.HandleFunc("/execute", worker.HandleExecute)

	slog.Info("worker listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
