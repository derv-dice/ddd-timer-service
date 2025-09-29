package service

import (
	"ddd-timer-service/config"
	cdrawer "ddd-timer-service/internal/pkg/calendar_drawer"
	"ddd-timer-service/internal/repository"

	"github.com/rs/zerolog"
)

type Service struct {
	conf           config.Config
	repo           repository.Repository
	logger         *zerolog.Logger
	calendarDrawer *cdrawer.CalendarDrawer
}

func New(repo repository.Repository, conf config.Config, logger *zerolog.Logger, cd *cdrawer.CalendarDrawer) *Service {
	return &Service{
		repo:           repo,
		conf:           conf,
		logger:         logger,
		calendarDrawer: cd,
	}
}

func (s *Service) Logger() *zerolog.Logger {
	return s.logger
}
