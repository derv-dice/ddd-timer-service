package app

import (
	"context"
	"ddd-timer-service/config"
	httpserver "ddd-timer-service/internal/api/http"
	"ddd-timer-service/internal/api/tg_bot"
	"ddd-timer-service/internal/repository"
	"ddd-timer-service/internal/service"
	"ddd-timer-service/internal/users_cache"
	"errors"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
)

func ConstructAndRun(ctx context.Context, wg *sync.WaitGroup, conf config.Config, logger *zerolog.Logger) error {
	// Подключение к БД
	logger.Info().Msg("init repository")
	sqliteRepo, err := repository.NewSQLiteRepository(conf.Database.Path, true)
	if err != nil {
		logger.Err(err).Msg("init repository failed")
		return err
	}

	memUsersCache := users_cache.NewImplUsersCacheMem()

	s := service.New(sqliteRepo, memUsersCache, conf, logger)

	logger.Info().Msg("loading users users_cache")
	if err = s.StartupUsersCache(ctx); err != nil {
		logger.Err(err).Msg("init users users_cache failed")
	}

	// Запуск http сервера
	logger.Info().Msg("start http server")
	httpServer := httpserver.NewImplServerGin(s)

	wg.Go(func() {
		errS := httpServer.Start(ctx, conf.Http.Addr)
		if errS != nil {
			logger.Err(errS).Msg("start http server failed")
		}
	})

	wg.Go(func() {
		<-ctx.Done()
		logger.Info().Msg("stop http server")

		errS := httpServer.Stop()
		if errS != nil && !errors.Is(errS, http.ErrServerClosed) {
			logger.Err(errS).Msg("stop http server failed")
			return
		}

		logger.Info().Msg("http server stopped")
	})

	// Запуск tg бота
	logger.Info().Msg("start telegram bot")
	tgBot := tg_bot.NewTelegramBot(s)

	wg.Go(func() {
		errS := tgBot.Start(ctx, conf.TGBot.Token)
		if errS != nil && !errors.Is(errS, context.Canceled) {
			logger.Err(errS).Msg("start telegram bot failed")
			return
		}
	})

	wg.Go(func() {
		<-ctx.Done()
		logger.Info().Msg("stop telegram bot")

		errS := tgBot.Stop()
		if errS != nil {
			logger.Err(errS).Msg("stop telegram bot failed")
			return
		}

		logger.Info().Msg("telegram bot stopped")
	})

	return nil
}
