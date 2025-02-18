package configs

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var onceInitRedisClient sync.Once

func InitRedisClient() {
	onceInitRedisClient.Do(func() {
		var tlsConfig *tls.Config
		if Env.RedisEnableSsl {
			tlsConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}

		rdb = redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%s:%s", Env.RedisHost, Env.RedisPort),
			Password:  Env.RedisPass,
			DB:        Env.RedisDB,
			TLSConfig: tlsConfig,
		})

		var ctx = context.Background()

		// Ping Redis to check connection
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatalf("failed to connect redis: %v", err)
		}
	})
}

func GetRedis() *redis.Client {
	if rdb == nil {
		log.Fatal("redis client not initialized")
	}
	return rdb
}
