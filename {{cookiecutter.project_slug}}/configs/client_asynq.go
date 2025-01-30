package configs

import (
	"crypto/tls"
	"fmt"

	"github.com/hibiken/asynq"
)

var AsynqClient *asynq.Client

func GetAsynqRedisOpt() asynq.RedisClientOpt {
	var tlsConfig *tls.Config
	if Env.RedisEnableSsl {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	return asynq.RedisClientOpt{
		Addr:      fmt.Sprintf("%s:%s", Env.RedisHost, Env.RedisPort),
		Password:  Env.RedisPass,
		TLSConfig: tlsConfig,
		DB:        Env.RedisTaskDB,
		//PoolSize:  1000,
	}
}

func InitAsynqClient() {
	AsynqClient = asynq.NewClient(GetAsynqRedisOpt())
}
