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

func ConstructAndRun(ctx context.Context, wg *sync.WaitGroup, conf config.Config, logger *zerolog.Logger) error {
	// Подключение к БД
	logger.Info().Msg("creating repository instance")
	sqliteRepo, err := repository.NewSQLiteRepository(ctx, conf.Database.Path, true)
	if err != nil {
		logger.Err(err).Msg("creating repository instance FAILED")
		return err
	}

	calendarDrawer := cdrawer.NewCalendarDrawer(ctx, conf.Limits.CalendarImg.MaxYears,
		conf.Limits.CalendarImg.CacheSizeMB)

	logger.Info().Msg("creating service instance")
	s := service.New(sqliteRepo, conf, logger, calendarDrawer)

	// Запуск http сервера
	logger.Info().Msg("starting http server")
	httpServer := httpserver.NewImplServerGin(s)

	wg.Go(func() {
		errS := httpServer.Start(ctx, conf.Http.Addr)
		if errS != nil {
			logger.Err(errS).Msg("start http server FAILED")
		}
	})

	wg.Go(func() {
		<-ctx.Done()
		logger.Info().Msg("stopping http server")

		errS := httpServer.Stop()
		if errS != nil && !errors.Is(errS, http.ErrServerClosed) {
			logger.Err(errS).Msg("stop http server FAILED")
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
			logger.Err(errS).Msg("start telegram bot FAILED")
			return
		}
	})

	wg.Go(func() {
		<-ctx.Done()
		logger.Info().Msg("stopping telegram bot")

		errS := tgBot.Stop()
		if errS != nil {
			logger.Err(errS).Msg("stopping telegram bot FAILED")
			return
		}

		logger.Info().Msg("telegram bot stopped")
	})

	return nil
}
