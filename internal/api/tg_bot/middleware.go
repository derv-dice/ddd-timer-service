package tg_bot

import (
	"context"
	"errors"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

func (i *implTelegramBot) skipNilMessagesMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *botmodels.Update) {
		if update == nil || update.Message == nil {
			err := errors.New("update.Message is nil")
			i.service.Logger().Err(err).Msg("nil message is received, skip it")
			return
		}

		next(ctx, bot, update)
	}
}

func (i *implTelegramBot) checkUserHasDatesMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		// Если длина текста равна regStrLen, пропустим сообщение. Это может быть попытка регистрации
		if len(update.Message.Text) == regStrLen {
			next(ctx, b, update)
			return
		}

		if !i.service.CheckUserHasServiceDates(ctx, update.Message.From.ID) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      mustRegisterMessage,
				ParseMode: botmodels.ParseModeMarkdown,
			})

			if err != nil {
				i.service.Logger().Err(err).Send()
			}

			return
		}

		next(ctx, b, update)
	}
}

func (i *implTelegramBot) logMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
		if update.Message != nil {
			i.service.Logger().Info().Int64("from", update.Message.From.ID).
				Str("msg", update.Message.Text).Msg("TGBOT Request")
		} else {
			i.service.Logger().Info().Int64("from", 0).Str("msg", update.Message.Text).
				Msg("TGBOT Request")
		}

		next(ctx, b, update)
	}
}
