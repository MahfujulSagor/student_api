package main

import (
	"context"
	"fmt"
	"github/mahfujulsagor/student_api/internal/config"
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
	logger.Init()

	//? Connect to DB

	//? Setup mux
	mux := http.NewServeMux()

	//? Setup routes
	mux.HandleFunc("POST /api/students", student.Create())

	//? Setup server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: mux,
	}
	logger.Info.Println("Server started on", server.Addr)

	//? Start server and listen for shutdown signal
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
