package config

import (
	"auth_service/internal/env"
	"errors"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/wt_2_pkg/logger"
)

type Config struct {
	Db     DbConfig
	Cache  CacheConfig
	Server ServerConfig
	Logger *logger.MyLogger
}

type DbConfig struct {
	PgUser         string `validate:"required"`
	PgPass         string `validate:"required"`
	PgHost         string `validate:"required"`
	PgPort         string `validate:"required"`
	PgDatabaseName string `validate:"required"`
	PgSSLMode      string `validate:"required"`
}

type CacheConfig struct {
	RedisUserName string `validate:"required"`
	RedisHost     string `validate:"required"`
	RedisPort     string `validate:"required"`
	RedisPass     string
	RedisDb       string `validate:"required"`
}

type ServerConfig struct {
	GrpcServerAddr string `validate:"required"`
	ServiceName    string `validate:"required"`
}

func LoadConfig() Config {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../../.env")
	}

	logger := logger.NewLogger()

	dbConfig := DbConfig{
		PgUser:         os.Getenv("PG_USER"),
		PgPass:         os.Getenv("PG_PASS"),
		PgHost:         os.Getenv("PG_HOST"),
		PgPort:         os.Getenv("PG_PORT"),
		PgDatabaseName: os.Getenv("PG_DATABASE_NAME"),
		PgSSLMode:      os.Getenv("PG_SSL_MODE"),
	}

	cacheConfig := CacheConfig{
		RedisUserName: os.Getenv("REDIS_USER_NAME"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPass:     os.Getenv("REDIS_PASS"),
		RedisDb:       os.Getenv("REDIS_DB"),
	}

	serverConfig := ServerConfig{
		GrpcServerAddr: os.Getenv("GRPC_SERVER_ADDR"),
		ServiceName:    os.Getenv("SERVICE_NAME"),
	}

	config := Config{
		Db:     dbConfig,
		Server: serverConfig,
		Cache:  cacheConfig,
		Logger: logger,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(config)
	if err != nil {
		var validatorErrs validator.ValidationErrors

		if errors.As(err, &validatorErrs) {
			log.Fatalf("config validation failed : %s", validatorErrs.Error())
		}
	}

	return config

}
