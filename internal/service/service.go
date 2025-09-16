package service

import (
	"ddd-timer-service/config"
	"ddd-timer-service/internal/repository"
	"ddd-timer-service/internal/users_cache"

	"github.com/rs/zerolog"
)

type Service struct {
	conf       config.Config
	repo       repository.Repository
	usersCache users_cache.UsersCache
	logger     *zerolog.Logger
}

func New(repo repository.Repository, usersCache users_cache.UsersCache, conf config.Config, logger *zerolog.Logger) *Service {
	return &Service{
		repo:       repo,
		conf:       conf,
		usersCache: usersCache,
		logger:     logger,
	}
}

func (s *Service) Logger() *zerolog.Logger {
	return s.logger
}
