package tg_bot

import (
	"context"
	"ddd-timer-service/internal/service"
	"ddd-timer-service/models"

	"github.com/go-telegram/bot"
)

type Bot interface {
	Start(ctx context.Context, token string) error
	Stop() error
}

type implTelegramBot struct {
	ok      bool
	pkgBot  *bot.Bot
	service *service.Service

	stop context.CancelFunc // bad, but necessary
}

func NewTelegramBot(service *service.Service) Bot {
	return &implTelegramBot{ok: true, service: service}
}

func (i *implTelegramBot) Start(ctx context.Context, token string) error {
	if i.ok == false {
		return models.ErrorNotInitialized
	}

	// Бот при создании уже начинает слать запросы в Telegram API,
	// поэтому этот код расположен здесь, а не в функции NewTelegramBot()
	options := []bot.Option{
		bot.WithDefaultHandler(i.accessLogMW(i.filterMW(i.defaultHandler))),
		bot.WithErrorsHandler(i.errorsHandler), // Просто заглушки
		bot.WithDebugHandler(i.debugHandler),   // Просто заглушки
	}

	var err error
	i.pkgBot, err = bot.New(token, options...)
	if nil != err {
		return err
	}

	i.pkgBot.RegisterHandler(bot.HandlerTypeMessageText, patternStart, bot.MatchTypeExact,
		i.accessLogMW(i.filterMW(i.startHandler)))

	i.pkgBot.RegisterHandler(bot.HandlerTypeMessageText, patternHelp, bot.MatchTypeExact,
		i.accessLogMW(i.filterMW(i.helpHandler)))

	i.pkgBot.RegisterHandler(bot.HandlerTypeMessageText, patternStats, bot.MatchTypeExact,
		i.accessLogMW(i.filterMW(i.statsHandler)))

	i.pkgBot.RegisterHandler(bot.HandlerTypeMessageText, patternUserInfo, bot.MatchTypeExact,
		i.accessLogMW(i.filterMW(i.getUserInfo)))

	i.pkgBot.RegisterHandler(bot.HandlerTypeMessageText, patternCells, bot.MatchTypeExact,
		i.accessLogMW(i.filterMW(i.cellsHandler)))

	i.pkgBot.RegisterHandler(bot.HandlerTypeMessageText, patternCalendar, bot.MatchTypeExact,
		i.accessLogMW(i.filterMW(i.calendarHandler)))

	i.pkgBot.RegisterHandler(bot.HandlerTypeMessageText, patternCalendarWithProgress, bot.MatchTypeExact,
		i.accessLogMW(i.filterMW(i.calendarWithProgressHandler)))

	var newCtx context.Context
	newCtx, i.stop = context.WithCancel(ctx)

	go i.pkgBot.Start(newCtx)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-newCtx.Done():
		return newCtx.Err()
	}
}

func (i *implTelegramBot) Stop() error {
	if i.ok == false {
		return models.ErrorNotInitialized
	}

	i.stop()

	return nil
}
