package configs

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var onceGetLogger sync.Once
var SentryClient *sentry.Client

func GetLogger() *zap.SugaredLogger {
	onceGetLogger.Do(func() {
		SentryClient = getSentryClient()
		logger = GetZapLogger(SentryClient)
	})
	return logger
}

func getZapSentryCore(client *sentry.Client) zapcore.Core {
	sentryCfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry, this level can't be higher than `Level`
		Tags: map[string]string{
			"component": "service",
		},
	}

	zapSentryCore, err := zapsentry.NewCore(sentryCfg, zapsentry.NewSentryClientFromClient(client))

	// don't use value if error was returned. Noop core will be replaced to nil soon.
	if err != nil {
		panic(err)
	}

	return zapSentryCore
}

func GetZapLogger(sentryClient *sentry.Client) *zap.SugaredLogger {

	cfg := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			TimeKey:      "time",
			LevelKey:     "level",
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeLevel:  CustomLevelEncoder,
			EncodeTime:   SyslogTimeEncoder,
		},
	}

	logger, _ := cfg.Build()

	if sentryClient != nil {
		zapSentryCore := getZapSentryCore(sentryClient)
		logger = zapsentry.AttachCoreToLogger(zapSentryCore, logger)
	}

	return logger.Sugar()
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func getSentryClient() *sentry.Client {
	if Env.SentryDSN == "" {
		return nil
	}

	log.Println("Init sentry: ")
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           Env.SentryDSN,
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("failed to init sentry: %v\n", err)
	}

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              Env.SentryDSN,
		TracesSampleRate: 1.0,
		Debug:            false,
	})

	if err != nil {
		log.Fatalf("failed to init sentry: %v", err)
	}

	return sentryClient
}
