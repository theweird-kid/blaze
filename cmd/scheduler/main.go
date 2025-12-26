package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/theweird-kid/blaze/internal/config"
	"github.com/theweird-kid/blaze/internal/db"
	"github.com/theweird-kid/blaze/internal/logger"
	"github.com/theweird-kid/blaze/internal/repository"
	"github.com/theweird-kid/blaze/internal/scheduler"
)

func main() {
	logger.Setup()

	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1) Connect to DB
	client, err := db.Connect(ctx, cfg.MongoURI)
	if err != nil {
		slog.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	// 2) Init Repos
	dbName := cfg.DBName
	jobRepo := repository.NewJobRepo(client.Database(dbName))
	jobRunRepo := repository.NewJobRunRepo(client.Database(dbName))

	// 3) Init Dispatcher
	workerURL := os.Getenv("WORKER_URL")
	if workerURL == "" {
		workerURL = "http://localhost:8080"
	}
	dispatcher := scheduler.NewHTTPDispatcher(workerURL)

	// 4) Start Scheduler
	s := scheduler.NewScheduler(jobRepo, jobRunRepo, dispatcher)

	go s.Start(ctx)

	// 5) Wait for shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	slog.Info("shutting down...")
}
