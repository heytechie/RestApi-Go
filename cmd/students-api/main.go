package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heytechie/Students-api-go/internal/config"
	"github.com/heytechie/Students-api-go/internal/http/handlers/student"
	"github.com/heytechie/Students-api-go/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	db, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("Error setting up database: %s", err)
	}
	slog.Info("Storage initialised", slog.String("storage_path", cfg.StoragePath))
	//setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(db))
	router.HandleFunc("GET /api/students/{id}", student.GetStudentById(db))
	router.HandleFunc("GET /api/students", student.GetAllStudents(db))
	//setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("Server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func(done chan os.Signal) {
		defer close(done)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}(done)

	<-done

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown successfully")

}
