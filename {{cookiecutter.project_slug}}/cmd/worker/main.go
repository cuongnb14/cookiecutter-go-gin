package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"{{ cookiecutter.project_slug }}/configs"

	"github.com/hibiken/asynq"
	"{{ cookiecutter.project_slug }}/internal/tasks"
)

func main() {
	configs.Bootstrap()

	var tlsConfig *tls.Config
	if configs.Env.RedisEnableSsl {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:      fmt.Sprintf("%s:%s", configs.Env.RedisHost, configs.Env.RedisPort),
			Password:  configs.Env.RedisPass,
			TLSConfig: tlsConfig,
			DB:        configs.Env.RedisTaskDB,
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 500,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				//"critical": 6,
				"default": 1,
				//"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.Use(tasks.LoggingMiddleware)
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
