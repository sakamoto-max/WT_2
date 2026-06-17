package config

import (
	// "email_service/internals/env"
	"os"
	"strconv"

	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type Config struct {
	Db       DbConfig
	Mq       MqConfig
	Consumer ConsumerConfig
	Logger   *logger.MyLogger
}

type DbConfig struct {
	UserName     string `validate:"required"`
	Pass         string `validate:"required"`
	Host         string `validate:"required"`
	Port         string `validate:"required"`
	DatabaseName string `validate:"required"`
	SSlMode      string `validate:"required"`
}

type MqConfig struct {
	UserName string `validate:"requied"`
	Pass     string `validate:"requied"`
	Host     string `validate:"requied"`
	Port     string `validate:"requied"`
}

type ConsumerConfig struct {
	NumberOfSenders int `validate:"required"`
	NumberOfWorkers int `validate:"required"`
}

func LoadConfig() Config {

	// stage := os.Getenv("STAGE")
	// if stage == "" {
	// 	env.Load("../.env")
	// }

	logger := logger.NewLogger()

	dbConfig := DbConfig{
		UserName:     os.Getenv("PG_USER"),
		Pass:         os.Getenv("PG_PASS"),
		Host:         os.Getenv("PG_HOST"),
		Port:         os.Getenv("PG_PORT"),
		DatabaseName: os.Getenv("PG_DATABASE_NAME"),
		SSlMode:      os.Getenv("PG_SSL_MODE"),
	}

	mqConfig := MqConfig{
		UserName: os.Getenv("MQ_USER_NAME"),
		Pass:     os.Getenv("MQ_PASS"),
		Host:     os.Getenv("MQ_HOSTNAME"),
		Port:     os.Getenv("MQ_PORT"),
	}

	numberOfSendersStr := os.Getenv("NUMBER_OF_SENDERS")
	numberOfSenders, err := strconv.Atoi(numberOfSendersStr)
	if err != nil {
		logger.Log.Fatalw("failed to get number of senders", zap.Error(err))
	}

	NumberOfWorkersStr := os.Getenv("NUMBER_OF_SENDERS")
	NumberOfWorkers, err := strconv.Atoi(NumberOfWorkersStr)
	if err != nil {
		logger.Log.Fatalw("failed to get number of workers", zap.Error(err))
	}

	consumerConfig := ConsumerConfig{
		NumberOfSenders: numberOfSenders,
		NumberOfWorkers: NumberOfWorkers,
	}

	config := Config{
		Db:       dbConfig,
		Mq:       mqConfig,
		Consumer: consumerConfig,
		Logger:   logger,
	}

	return config

}
