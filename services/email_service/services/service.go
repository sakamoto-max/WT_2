package services

import (
	
	"github.com/sakamoto-max/wt_2-pkg/logger"


	"go.uber.org/zap"
)

type Service struct {
	logger *logger.MyLogger
}

func NewService(logger *logger.MyLogger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) SendWelcomeEmail(email string) error {
	s.logger.Log.Infow(
		"email sent",
		 zap.String("email", email),
	)

	return nil
}