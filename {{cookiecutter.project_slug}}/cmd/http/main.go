package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/server"
)

func main() {
	configs.Bootstrap()
	logger := configs.GetLogger()

	configs.InitAsynqClient()
	defer configs.AsynqClient.Close()

	var s = server.Server{}
	s.Initialize()

	go func() {
		s.Run()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 30 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	if configs.SentryClient != nil {
		configs.SentryClient.Flush(2 * time.Second)
	}

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %s", err)
	}

	logger.Info("Server exiting")
}
