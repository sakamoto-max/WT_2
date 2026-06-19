package config

import (
	"errors"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	// "github.com/sakamoto-max/wt_2/api_gateway/internals/env"
	"github.com/sakamoto-max/wt_2_pkg/logger"
)

type Config struct {
	Logger              *logger.MyLogger
	HttpServer          HttpConfig
	AuthServerConfig    ClientConfig
	PlanServerConfig    ClientConfig
	ExerServerConfig    ClientConfig
	TrackerServerConfig ClientConfig
	AuthConfig          AuthConfig
}

type HttpConfig struct {
	Addr       string `validate:"required"`
	ServerName string `validate:"required"`
}

type ClientConfig struct {
	Host string `validate:"required"`
	Addr string `validate:"required"`
}

type AuthConfig struct {
	SecretKey string `validate:"required"`
}

func LoadConfig() Config {

	// env.Load("../../.env")

	logger := logger.NewLogger()

	httpServerConfig := HttpConfig{
		Addr:       os.Getenv("HTTP_SERVER_ADDR"),
		ServerName: os.Getenv("SERVICE_NAME"),
	}

	authServerConfig := ClientConfig{
		Host: os.Getenv("AUTH_GRPC_CLIENT_HOST"),
		Addr: os.Getenv("AUTH_GRPC_CLIENT_ADDR"),
	}

	planServerConfig := ClientConfig{
		Host: os.Getenv("PLAN_GRPC_CLIENT_HOST"),
		Addr: os.Getenv("PLAN_GRPC_CLIENT_ADDR"),
	}
	exerServerConfig := ClientConfig{
		Host: os.Getenv("EXER_GRPC_CLIENT_HOST"),
		Addr: os.Getenv("EXER_GRPC_CLIENT_ADDR"),
	}
	trackerServerConfig := ClientConfig{
		Host: os.Getenv("TRACKER_GRPC_CLIENT_HOST"),
		Addr: os.Getenv("TRACKER_GRPC_CLIENT_ADDR"),
	}

	authConfig := AuthConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	config := Config{
		HttpServer:          httpServerConfig,
		AuthServerConfig:    authServerConfig,
		PlanServerConfig:    planServerConfig,
		TrackerServerConfig: trackerServerConfig,
		ExerServerConfig:    exerServerConfig,
		AuthConfig:          authConfig,
		Logger:              logger,
	}

	validate := validator.New()

	err := validate.Struct(config)
	if err != nil {
		var validatorErrs validator.ValidationErrors

		if errors.As(err, &validatorErrs) {
			log.Fatalf("config validation failed : %s", validatorErrs.Error())
		}
	}

	return config
}
