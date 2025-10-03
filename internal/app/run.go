package app

import (
	"context"
	"ddd-timer-service/config"
	httpserver "ddd-timer-service/internal/api/http"
	"ddd-timer-service/internal/api/tg_bot"
	cdrawer "ddd-timer-service/internal/pkg/calendar_drawer"
	"ddd-timer-service/internal/repository"
	"ddd-timer-service/internal/service"
	"errors"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
)

func Run(ctx context.Context, wg *sync.WaitGroup, conf config.Config, logger *zerolog.Logger) error {
	// Подключение к БД
	logger.Info().Msg("creating repository instance")

	sqliteRepo, err := repository.NewSQLiteRepository(ctx, conf.Database.Path, true)
	if err != nil {
		logger.Error().Err(err).Msg("repository instance failed")
		return err
	}

	calendarDrawer := cdrawer.NewCalendarDrawer(ctx, conf.Limits.CalendarImg.MaxYears,
		conf.Limits.CalendarImg.CacheSizeMB)

	logger.Info().Msg("creating service instance")
	s := service.New(sqliteRepo, conf, calendarDrawer)

	// Запуск http сервера
	logger.Info().Msg("starting http server")
	httpServer := httpserver.NewImplServerGin(s)

	wg.Go(func() {
		errS := httpServer.Start(ctx, conf.Http.Addr)
		if errS != nil {
			logger.Error().Err(errS).Msg("starting http server failed")
		}
	})

	wg.Go(func() {
		<-ctx.Done()
		logger.Info().Msg("stopping http server")

		errS := httpServer.Stop()
		if errS != nil && !errors.Is(errS, http.ErrServerClosed) {
			logger.Error().Err(errS).Msg("stopping http server failed")

			return
		}

		logger.Info().Msg("http server stopped")
	})

	// Запуск tg бота
	logger.Info().Msg("starting telegram bot")
	tgBot := tg_bot.NewTelegramBot(s)

	wg.Go(func() {
		errS := tgBot.Start(ctx, conf.TGBot.Token)
		if errS != nil && !errors.Is(errS, context.Canceled) {
			logger.Error().Err(errS).Msg("starting telegram bot failed")

			return
		}
	})

	wg.Go(func() {
		<-ctx.Done()
		logger.Info().Msg("stopping telegram bot")

		errS := tgBot.Stop()
		if errS != nil {
			logger.Error().Err(errS).Msg("stopping telegram bot failed")

			return
		}

		logger.Info().Msg("telegram bot stopped")
	})

	return nil
}
