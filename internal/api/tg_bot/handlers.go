package tg_bot

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

func (i *implTelegramBot) startHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	userID := update.Message.From.ID

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Привет, *%s*", bot.EscapeMarkdown(update.Message.From.FirstName)),
		ParseMode: botmodels.ParseModeMarkdown,
	})

	if !i.service.CheckUserHasServiceDates(ctx, userID) {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      mustRegisterMessage,
			ParseMode: botmodels.ParseModeMarkdownV1,
		})
	}
}

func (i *implTelegramBot) defaultHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	userID := update.Message.From.ID

	// Более простая проверка, чем на каждое сообщени вызывать regexp
	if len(update.Message.Text) == regStrLen {
		if datesRegex.MatchString(update.Message.Text) {
			dates := strings.Split(update.Message.Text, " ")

			if err := i.service.SetUserDatesFromStringMessage(ctx, userID, dates[0], dates[1]); err != nil {
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   fmt.Sprintf("Не удалось сохранить даты, ошибка: %s", err.Error()),
				})

				return
			}

			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Даты начала и окончания службы изменены",
			})
		}
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Чтобы узнать статистику используй команду /stats",
		ParseMode: botmodels.ParseModeMarkdown,
	})
}

func (i *implTelegramBot) statsHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	userID := update.Message.From.ID

	s, err := i.service.GetUserStats(ctx, userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   s.PrettyShort(),
	})
}

func (i *implTelegramBot) getUserInfo(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	userID := update.Message.From.ID

	user, err := i.service.GetUser(ctx, userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   user.String(),
	})
}

func (i *implTelegramBot) helpHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      helpMessage1,
		ParseMode: botmodels.ParseModeMarkdown,
	})
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpMessage,
	})
}

func (i *implTelegramBot) cellsHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	userID := update.Message.From.ID

	img, err := i.service.GenerateCellsPNG(ctx, userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		return
	}

	stats, err := i.service.GetUserStats(context.Background(), userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		return
	}

	media := &botmodels.InputMediaPhoto{
		Media: "attach://cells.png",
		Caption: fmt.Sprintf("Дней прошло %d из %d, сталось: %d",
			int(math.Floor(stats.PassedDays())), stats.TotalDays(), int(math.Ceil(stats.LeftDays()))),
		MediaAttachment: bytes.NewReader(img),
	}

	_, _ = b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: update.Message.Chat.ID,
		Media: []botmodels.InputMedia{
			media,
		},
	})
}

func (i *implTelegramBot) calendarHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	userID := update.Message.From.ID

	img, err := i.service.GenerateCalendarPNG(ctx, userID, false)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		return
	}

	media := &botmodels.InputMediaPhoto{
		Media:           "attach://ddd_calendar.png",
		Caption:         "Посезонный календарь на весь период службы. Можешь распечатать его и отмечать дни вручную, либо вызвать команду /calendar_with_progress и получить картунку уже заполненного.",
		MediaAttachment: bytes.NewReader(img),
	}

	_, _ = b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: update.Message.Chat.ID,
		Media: []botmodels.InputMedia{
			media,
		},
	})
}

func (i *implTelegramBot) calendarWithProgressHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	userID := update.Message.From.ID

	img, err := i.service.GenerateCalendarPNG(ctx, userID, true)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		return
	}

	media := &botmodels.InputMediaPhoto{
		Media:           "attach://ddd_calendar_with_progress.png",
		Caption:         "Посезонный календарь на весь период службы с отметками о прошедших днях. Чтобы получить такой же, но без отметок, вызови команду /calendar",
		MediaAttachment: bytes.NewReader(img),
	}

	_, _ = b.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: update.Message.Chat.ID,
		Media: []botmodels.InputMedia{
			media,
		},
	})
}

func (i *implTelegramBot) errorsHandler(err error) {
	i.service.Logger().Err(err).Msg("TGBOT errorsHandler")
}

func (i *implTelegramBot) debugHandler(format string, args ...any) {
	i.service.Logger().Debug().Str("debug_msg", fmt.Sprintf(format, args...)).Msg("TGBOT debugHandler")
}
