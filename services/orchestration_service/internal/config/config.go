package config

import (
	"errors"
	"log"
	"orchestration_service/internal/env"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type Config struct {
	Dbs      Dbs
	Mq       MqConfig
	Logger   *logger.MyLogger
	Consumer ConsumerConfig
	Cache    Cache
}

type MqConfig struct {
	MqUserName string `validate:"required"`
	MqPass     string `validate:"required"`
	MqHostName string `validate:"required"`
	MqPort     string `validate:"required"`
}

type Dbs struct {
	Auth    DbConfig
	Tracker DbConfig
}

type DbConfig struct {
	PgUser         string `validate:"required"`
	PgPass         string `validate:"required"`
	PgHost         string `validate:"required"`
	PgPort         string `validate:"required"`
	PgDatabaseName string `validate:"required"`
	PgSSLMode      string `validate:"required"`
}

type Cache struct {
	UserName string `validate:"required"`
	Password string `validate:"required"`
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Db       string `validate:"required"`
}

type ConsumerConfig struct {
	NumberOfSenders int `validate:"required"`
	NumberOfWorkers int `validate:"required"`
}

func LoadConfig() Config {

	stage := os.Getenv("STAGE")
	if stage == "" {
		env.Load("../.env")
	}

	logger := logger.NewLogger()

	authDbConfig := DbConfig{
		PgUser:         os.Getenv("AUTH_USERNAME"),
		PgPass:         os.Getenv("AUTH_PASSWORD"),
		PgHost:         os.Getenv("AUTH_HOSTNAME"),
		PgPort:         os.Getenv("AUTH_PORT"),
		PgDatabaseName: os.Getenv("AUTH_DATABASE_NAME"),
		PgSSLMode:      os.Getenv("AUTH_SSL_MODE"),
	}

	trackerDbConfig := DbConfig{
		PgUser:         os.Getenv("TRACKER_USERNAME"),
		PgPass:         os.Getenv("TRACKER_PASSWORD"),
		PgHost:         os.Getenv("TRACKER_HOSTNAME"),
		PgPort:         os.Getenv("TRACKER_PORT"),
		PgDatabaseName: os.Getenv("TRACKER_DATABASE_NAME"),
		PgSSLMode:      os.Getenv("TRACKER_SSL_MODE"),
	}

	dbs := Dbs{
		Auth:    authDbConfig,
		Tracker: trackerDbConfig,
	}

	mqConfig := MqConfig{
		MqUserName: os.Getenv("MQ_USER_NAME"),
		MqPass:     os.Getenv("MQ_PASS"),
		MqHostName: os.Getenv("MQ_HOSTNAME"),
		MqPort:     os.Getenv("MQ_PORT"),
	}

	cache := Cache{
		UserName: os.Getenv("REDIS_USER_NAME"),
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
		Db:       os.Getenv("REDIS_DB"),
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

	consumer := ConsumerConfig{
		NumberOfSenders: numberOfSenders,
		NumberOfWorkers: numberOfWorkers,
	}

	config := Config{
		Logger:   logger,
		Dbs:      dbs,
		Consumer: consumer,
		Mq:       mqConfig,
		Cache:    cache,
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
