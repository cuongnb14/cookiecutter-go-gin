package configs

import (
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/phsym/console-slog"
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
)

var logger *slog.Logger
var onceInitLogger sync.Once

func InitLogger() {
	onceInitLogger.Do(func() {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              Env.SentryDSN,
			EnableTracing:    false,
			TracesSampleRate: 1.0,
			Environment:      Env.Stage,
		})
		if err != nil {
			log.Fatal(err)
		}

		logger = slog.New(
			slogmulti.Fanout(
				console.NewHandler(os.Stdout, &console.HandlerOptions{Level: slog.LevelInfo, NoColor: true}),
				slogsentry.Option{Level: slog.LevelWarn}.NewSentryHandler(),
			))
		slog.SetDefault(logger)

	})
}

func GetLogger() *slog.Logger {
	if logger == nil {
		log.Fatal("logger not initialized")
	}
	return logger
}
