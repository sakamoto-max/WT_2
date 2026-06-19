package config

import (
	"errors"
	"log"
	"os"
	// "plan_service/internal/env"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type Config struct {
	Db            DbConfig
	Cache         CacheConfig
	Server        ServerConfig
	Mq            MqConfig
	OtherServices OtherServices
	Logger        *logger.MyLogger
	Consumer
}

type Consumer struct {
	NumberOfSenders int `validate:"required"`
	NumberOfWorkers int `validate:"required"`
}

type MqConfig struct {
	MqUserName string `validate:"required"`
	MqPass     string `validate:"required"`
	MqHostName string `validate:"required"`
	MqPort     string `validate:"required"`
}

type OtherServices struct {
	ExerServiceHost string `validate:"required"`
	ExerServiceAddr string `validate:"required"`
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

	// env.Load("../../.env")

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

	mqConfig := MqConfig{
		MqUserName: os.Getenv("MQ_USER_NAME"),
		MqPass:     os.Getenv("MQ_PASS"),
		MqHostName: os.Getenv("MQ_HOSTNAME"),
		MqPort:     os.Getenv("MQ_PORT"),
	}

	otherServices := OtherServices{
		ExerServiceHost: os.Getenv("EXERCISE_SERVICE_HOST"),
		ExerServiceAddr: os.Getenv("EXERCISE_SERVICE_ADDR"),
	}

	numberOfSendersStr := os.Getenv("NUMBER_OF_SENDERS")
	numberOfSenders, err := strconv.Atoi(numberOfSendersStr)
	if err != nil {
		logger.Log.Fatalw("failed to get number of senders", zap.Error(err))
	}

	numberOfWorkersStr := os.Getenv("NUMBER_OF_WORKERS")
	numberOfWorkers, err := strconv.Atoi(numberOfWorkersStr)
	if err != nil {
		logger.Log.Fatalw("failed to get number of workers", zap.Error(err))
	}

	consumer := Consumer{
		NumberOfSenders: numberOfSenders,
		NumberOfWorkers: numberOfWorkers,
	}

	config := Config{
		Db:            dbConfig,
		Server:        serverConfig,
		Cache:         cacheConfig,
		Logger:        logger,
		Mq:            mqConfig,
		OtherServices: otherServices,
		Consumer:      consumer,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		var validatorErrs validator.ValidationErrors

		if errors.As(err, &validatorErrs) {
			log.Fatalf("config validation failed : %s", validatorErrs.Error())
		}
	}

	return config

}
