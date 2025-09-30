package tg_bot

import (
	"context"
	"ddd-timer-service/internal/pkg/tracelog"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

// MW = Middleware
//
// rootMiddleware -> traceIDMW -> accessLogMW -> skipNilMW -> handler
func (i *implTelegramBot) rootMiddleware(next bot.HandlerFunc, pass ...bool) bot.HandlerFunc {
	if len(pass) > 0 {
		return i.traceIDMW(i.accessLogMW(next))
	}

	return i.traceIDMW(i.accessLogMW(i.checkUserHasDatesMW(next)))
}

func (i *implTelegramBot) traceIDMW(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *botmodels.Update) {
		tl, ctx := tracelog.Begin(ctx, "tgbot.root")
		defer tl.End()

		next(ctx, bot, update)
	}
}

func (i *implTelegramBot) checkUserHasDatesMW(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		// Если длина текста равна regStrLen, пропустим сообщение. Это может быть попытка регистрации
		if len(update.Message.Text) == regStrLen {
			next(ctx, b, update)
			return
		}

		if !i.service.CheckUserHasServiceDates(ctx, update.Message.From.ID) {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      mustRegisterMessage,
				ParseMode: botmodels.ParseModeMarkdown,
			})

			return
		}

		next(ctx, b, update)
	}
}

func (i *implTelegramBot) accessLogMW(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		tl, ctx := tracelog.Begin(ctx, "tgbot.accessLog")
		defer tl.End()

		if update.Message != nil {

			tl.Trace("new message", tracelog.Int(logKeyUserID, int(update.Message.From.ID)),
				tracelog.String(logKeyMessage, update.Message.Text))
		} else {
			tl.Warn("non message update")
			return
		}

		next(ctx, b, update)
	}
}
