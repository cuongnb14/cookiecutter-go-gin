package configs

import (
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/hibiken/asynq"
)

var asynqClient *asynq.Client
var onceInitAsynqClient sync.Once

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
	onceInitAsynqClient.Do(func() {
		asynqClient = asynq.NewClient(GetAsynqRedisOpt())
	})
}

func GetAsynqClient() *asynq.Client {
	if asynqClient == nil {
		log.Fatal("asynq client not initialized")
	}
	return asynqClient
}
