package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Swiffy12/taskify/internal/config"
	taskhandler "github.com/Swiffy12/taskify/internal/http-server/handlers/task"
	taskservice "github.com/Swiffy12/taskify/internal/http-server/services/task"
	"github.com/Swiffy12/taskify/internal/lib/logger/sl"
	"github.com/Swiffy12/taskify/internal/storage/postgresql"
	"github.com/go-chi/chi"
)

func main() {
	config := config.MustLoad()

	log := setupLogger(config.Env)

	storage, err := postgresql.New(
		config.Storage.DBHost,
		config.Storage.DBPort,
		config.Storage.DBUser,
		config.Storage.DBPassword,
		config.Storage.DBName,
	)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer storage.DB.Close()

	taskService := taskservice.New(storage)
	taskHandler := taskhandler.New(log, taskService)
	router := chi.NewRouter()

	router.Post("/tasks", taskHandler.CreateTask)
	router.Get("/tasks/{id}", taskHandler.GetTask)
	router.Get("/tasks", taskHandler.GetAllTasks)
	router.Patch("/tasks/{id}", taskHandler.UpdateTask)
	router.Delete("/tasks/{id}", taskHandler.DeleteTask)

	log.Info("starting server")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Handler:      router,
		Addr:         config.Host + ":" + config.Port,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
