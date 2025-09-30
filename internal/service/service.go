package service

import (
	"ddd-timer-service/config"
	cdrawer "ddd-timer-service/internal/pkg/calendar_drawer"
	"ddd-timer-service/internal/repository"
)

type Service struct {
	conf           config.Config
	repo           repository.Repository
	calendarDrawer *cdrawer.CalendarDrawer
}

func New(repo repository.Repository, conf config.Config, cd *cdrawer.CalendarDrawer) *Service {
	return &Service{
		repo:           repo,
		conf:           conf,
		calendarDrawer: cd,
	}
}
