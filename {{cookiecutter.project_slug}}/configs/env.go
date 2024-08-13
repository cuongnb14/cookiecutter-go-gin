package configs

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	Stage   string `mapstructure:"STAGE"`
	GinMode string `mapstructure:"GIN_MODE"`
	Port    string `mapstructure:"PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPass     string `mapstructure:"DB_PASS"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSslMode  string `mapstructure:"DB_SSL_MODE"`
	DBLogLevel int    `mapstructure:"DB_LOG_LEVEL"`

	SentryDSN              string `mapstructure:"SENTRY_DSN"`
	EnableLogRequestDetail bool   `mapstructure:"ENABLE_LOG_REQUEST_DETAIL"`

	JwtSecret string `mapstructure:"JWT_SECRET"`

	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPort      string `mapstructure:"REDIS_PORT"`
	RedisPass      string `mapstructure:"REDIS_PASS"`
	RedisDB        int    `mapstructure:"REDIS_DB"`
	RedisEnableSsl bool   `mapstructure:"REDIS_ENABLE_SSL"`
	RedisTaskDB    int    `mapstructure:"REDIS_TASK_DB"`
}

var Env = EnvConfig{}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filename))
	stage := os.Getenv("STAGE")

	if strings.HasSuffix(os.Args[0], ".test") {
		stage = "testing"
	}

	envFile := rootDir + "/env/" + stage + ".env"

	log.Printf("Use env file: %s", envFile)
	viper.SetConfigFile(envFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	err = viper.Unmarshal(&Env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}
}
