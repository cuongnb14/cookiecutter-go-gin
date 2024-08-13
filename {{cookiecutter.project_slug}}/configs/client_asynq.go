package configs

import (
	"crypto/tls"
	"fmt"
	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client

func InitAsynqClient() {
	var tlsConfig *tls.Config
	if Env.RedisEnableSsl {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:      fmt.Sprintf("%s:%s", Env.RedisHost, Env.RedisPort),
		Password:  Env.RedisPass,
		TLSConfig: tlsConfig,
		DB:        Env.RedisTaskDB,
	})

}
