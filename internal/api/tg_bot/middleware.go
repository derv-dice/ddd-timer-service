package tg_bot

import (
	"context"
	"ddd-timer-service/internal/pkg/tracelog"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

func (i *implTelegramBot) accessLogMW(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		tl, ctx := tracelog.Begin(ctx, "TGBOT")
		defer tl.End()

		var msg string

		if update.Message != nil {

			msg = update.Message.Text

			tl.AddAttributes(tracelog.Int(logKeyUserID, int(update.Message.From.ID)),
				tracelog.Int(logKeyChatID, int(update.Message.Chat.ID)))

		} else {
			tl.Warn("nil message update")
			return
		}

		next(ctx, b, update)
		tl.InfoWithDuration("message processed", tracelog.String(logKeyMessage, msg))
	}
}

func (i *implTelegramBot) filterMW(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		// Если длина текста равна regStrLen, пропустим сообщение. Это может быть попытка регистрации
		if len(update.Message.Text) == regStrLen {
			next(ctx, b, update)
			return
		}

		// Проверка, что пользователь зарегистрирован
		if update.Message != nil && filterCheckUserHasServiceDates[update.Message.Text] {
			if !i.service.CheckUserHasServiceDates(ctx, update.Message.From.ID) {
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:    update.Message.Chat.ID,
					Text:      mustRegisterMessage,
					ParseMode: botmodels.ParseModeMarkdown,
				})

				return
			}
		}

		next(ctx, b, update)
	}
}

var filterCheckUserHasServiceDates = map[string]bool{
	patternStart:                true,
	patternHelp:                 false,
	patternStats:                false,
	patternUserInfo:             false,
	patternCells:                false,
	patternCalendar:             false,
	patternCalendarWithProgress: false,
}
