package configs

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

var rdb *redis.Client
var onceGetRedis sync.Once

func GetRedis() *redis.Client {
	onceGetRedis.Do(func() {
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

	return rdb
}
