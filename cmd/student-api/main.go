package main

import (
	"context"
	"fmt"
	"github/mahfujulsagor/student_api/internal/config"
	"github/mahfujulsagor/student_api/internal/db/sqlite"
	"github/mahfujulsagor/student_api/internal/http/handlers/student"
	"github/mahfujulsagor/student_api/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//? Load config
	cfg := config.MustLoad()

	//? Initialize logger
	logger.Init(cfg, cfg.Env)

	//? Setup database
	db, err := sqlite.New(cfg)
	if err != nil {
		logger.Error.Fatal("Failed to connect to database:", err)
		return
	}
	logger.Info.Println("Connected to database", "env:", cfg.Env)

	//? Setup mux
	mux := http.NewServeMux()

	//? Setup routes
	mux.HandleFunc("POST /api/students", student.New(db))
	mux.HandleFunc("GET /api/students/{id}", student.GetByID(db))
	mux.HandleFunc("GET /api/students", student.GetList(db))
	mux.HandleFunc("PUT /api/students/{id}", student.UpdateByID(db))

	//? Setup server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: mux,
	}

	//? Start server and listen for shutdown signal
	logger.Info.Println("Server started on", server.Addr)
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error.Fatal("Failed to start server:", err)
		}
	}()
	<-done

	logger.Info.Println("Server shutting down...")

	//? Shutdown server gracefully within 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error.Fatal("Server forced to shutdown:", err)
	}
	logger.Info.Println("Server shut down gracefully")
}
