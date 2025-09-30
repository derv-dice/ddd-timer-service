package tg_bot

import (
	"bytes"
	"context"
	"ddd-timer-service/internal/pkg/tracelog"
	"fmt"
	"math"
	"strings"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

func (i *implTelegramBot) startHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "tgbot/startHandler")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)))

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
	tl, ctx := tracelog.Begin(ctx, "tgbot/defaultHandler")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)), tracelog.String(logKeyMessage, update.Message.Text))

	// Более простая проверка, чем на каждое сообщении вызывать regexp
	if len(update.Message.Text) == regStrLen {
		if datesRegex.MatchString(update.Message.Text) {
			dates := strings.Split(update.Message.Text, " ")

			if err := i.service.SetUserDatesFromStringMessage(ctx, userID, dates[0], dates[1]); err != nil {
				_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   fmt.Sprintf("Не удалось сохранить даты, ошибка: %s", err.Error()),
				})

				tl.AddError(err)
				return
			}

			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Даты начала и окончания службы изменены",
			})

			tl.Info("user dates has been changed")
		}
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Чтобы узнать статистику используй команду /stats",
		ParseMode: botmodels.ParseModeMarkdown,
	})
}

func (i *implTelegramBot) statsHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "tgbot/statsHandler")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)))

	s, err := i.service.GetUserStats(ctx, userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		tl.AddError(err)
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   s.PrettyShort(),
	})
}

func (i *implTelegramBot) getUserInfo(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "tgbot/getUserInfo")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)))

	user, err := i.service.GetUser(ctx, userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		tl.AddError(err)
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   user.String(),
	})
}

func (i *implTelegramBot) helpHandler(ctx context.Context, b *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "tgbot/helpHandler")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)))

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
	tl, ctx := tracelog.Begin(ctx, "tgbot/cellsHandler")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)))

	img, err := i.service.GenerateCellsPNG(ctx, userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		tl.AddError(err)
		return
	}

	stats, err := i.service.GetUserStats(context.Background(), userID)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		tl.AddError(err)
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
	tl, ctx := tracelog.Begin(ctx, "tgbot/calendarHandler")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)))

	img, err := i.service.GenerateCalendarPNG(ctx, userID, false)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Ошибка: %s", err.Error()),
		})

		tl.AddError(err)
		return
	}

	media := &botmodels.InputMediaPhoto{
		Media:           "attach://ddd_calendar.png",
		Caption:         "Посезонный календарь на весь период службы. Можешь распечатать его и отмечать дни вручную, либо вызвать команду /calendar_with_progress и получить картинку уже заполненного календаря.",
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
	tl, ctx := tracelog.Begin(ctx, "tgbot/calendarWithProgressHandler")
	defer tl.End()

	userID := update.Message.From.ID
	tl.AddAttributes(tracelog.Int(logKeyUserID, int(userID)))

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

// errorsHandler - заглушка для bot.ErrorsHandler
func (i *implTelegramBot) errorsHandler(err error) {
	if strings.Contains(err.Error(), "getUpdates\": context canceled") {
		return
	}

	tl, _ := tracelog.Begin(context.Background(), "tgbot.errorsHandler")
	tl.AddError(err)
	tl.End()
}

// debugHandler - заглушка для bot.DebugHandler
func (i *implTelegramBot) debugHandler(_ string, _ ...any) {
	tl, _ := tracelog.Begin(context.Background(), "tgbot.debugHandler")
	tl.End()
}
