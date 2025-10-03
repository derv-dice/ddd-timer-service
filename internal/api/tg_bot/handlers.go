package tg_bot

import (
	"context"
	"ddd-timer-service/internal/pkg/tracelog"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/rs/zerolog"
)

func (i *implTelegramBot) startHandler(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/startHandler")
	defer tl.End()

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	if err := i.SendMessage(ctx, chatID, pmMDv2, newHelloMessage(update.Message.From.FirstName)); err != nil {
		tl.AddError(err)
		return
	}

	if !i.service.CheckUserHasServiceDates(ctx, userID) {
		if err := i.SendMessage(ctx, chatID, pmMDv1, mustRegisterMessage); err != nil {
			tl.AddError(err)
			return
		}

		tl.InfoWithDuration("user without service dates")
		return
	}

	if err := i.SendMessage(ctx, chatID, pmMDv2, basicUsageMessage); err != nil {
		tl.AddError(err)
		return
	}
}

func (i *implTelegramBot) defaultHandler(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/defaultHandler")
	defer tl.End()

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Проверка, что пользователь зарегистрирован
	if update.Message != nil && filterCheckUserHasServiceDates[update.Message.Text] {
		if !i.service.CheckUserHasServiceDates(ctx, update.Message.From.ID) {
			if err := i.SendMessage(ctx, chatID, pmMDv1, mustRegisterMessage); err != nil {
				tl.AddError(err)
				return
			}

			tl.InfoWithDuration("user without service dates")
			return
		}
	}

	// Проверяем, возможно пользователь хочет зарегистрироваться или поменять даты
	if len(update.Message.Text) == regStrLen {
		if datesRegex.MatchString(update.Message.Text) {
			dates := strings.Split(update.Message.Text, " ")
			from, to, err := stringDatesToTime(dates[0], dates[1])
			if err != nil {
				tl.AddError(err, zerolog.WarnLevel)
				return
			}

			if !from.Before(to) || from.Round(time.Hour*24).Equal(to.Round(time.Hour*24)) {
				err = ErrBadDates
				tl.AddError(err, zerolog.WarnLevel)
				return
			}

			if err = i.service.SaveUser(ctx, userID, from, to); err != nil {
				tl.AddError(err)

				if sErr := i.SendErrorMessage(ctx, chatID, err); sErr != nil {
					tl.AddError(sErr)
					return
				}

				return
			}

			tl.InfoWithDuration("user dates has been changed")

			if sErr := i.SendMessage(ctx, chatID, pmNone, "даты начала и окончания службы успешно изменены"); sErr != nil {
				tl.AddError(sErr)
				return
			}
		}
	}

	if sErr := i.SendMessage(ctx, chatID, pmMDv2, basicUsageMessage); sErr != nil {
		tl.AddError(sErr)
		return
	}
}

func (i *implTelegramBot) statsHandler(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/statsHandler")
	defer tl.End()

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	s, err := i.service.GetUserStats(ctx, userID)
	if err != nil {
		tl.AddError(err)

		if sErr := i.SendErrorMessage(ctx, chatID, err); sErr != nil {
			tl.AddError(sErr)
			return
		}

		return
	}

	if sErr := i.SendMessage(ctx, chatID, pmNone, s.PrettyShort()); sErr != nil {
		tl.AddError(sErr)
		return
	}
}

func (i *implTelegramBot) getUserInfo(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/getUserInfo")
	defer tl.End()

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	user, err := i.service.GetUser(ctx, userID)
	if err != nil {
		tl.AddError(err)

		if sErr := i.SendErrorMessage(ctx, chatID, err); sErr != nil {
			tl.AddError(sErr)
			return
		}

		return
	}

	if sErr := i.SendMessage(ctx, chatID, pmNone, user.String()); sErr != nil {
		tl.AddError(sErr)
		return
	}
}

func (i *implTelegramBot) helpHandler(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/helpHandler")
	defer tl.End()

	chatID := update.Message.Chat.ID

	if sErr := i.SendMessage(ctx, chatID, pmMDv2, helpMessage); sErr != nil {
		tl.AddError(sErr)
		return
	}
}

func (i *implTelegramBot) cellsHandler(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/cellsHandler")
	defer tl.End()

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	stats, err := i.service.GetUserStats(ctx, userID)
	if err != nil {
		tl.AddError(err)

		if sErr := i.SendErrorMessage(ctx, chatID, err); sErr != nil {
			tl.AddError(sErr)
			return
		}

		return
	}

	img, err := i.service.GenerateCellsPNG(ctx, userID)
	if err != nil {
		tl.AddError(err)

		if sErr := i.SendErrorMessage(ctx, chatID, err); sErr != nil {
			tl.AddError(sErr)
			return
		}

		return
	}

	caption := fmt.Sprintf("Дней прошло %d из %d, сталось: %d", int(math.Floor(stats.PassedDays())),
		stats.TotalDays(), int(math.Ceil(stats.LeftDays())))

	name := fmt.Sprintf("cells_%s", time.Now().Format("2006-01-02"))

	if sErr := i.SendImagePNG(ctx, chatID, name, caption, img); sErr != nil {
		tl.AddError(sErr)
		return
	}
}

func (i *implTelegramBot) calendarHandler(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/calendarHandler")
	defer tl.End()

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	img, err := i.service.GenerateCalendarPNG(ctx, userID, false)
	if err != nil {
		tl.AddError(err)

		if sErr := i.SendErrorMessage(ctx, chatID, err); sErr != nil {
			tl.AddError(sErr)
			return
		}

		return
	}

	caption := captionCalendar
	name := fmt.Sprintf("calendar_%s", time.Now().Format("2006-01-02"))

	if sErr := i.SendImagePNG(ctx, chatID, name, caption, img); sErr != nil {
		tl.AddError(sErr)
		return
	}
}

func (i *implTelegramBot) calendarWithProgressHandler(ctx context.Context, _ *bot.Bot, update *botmodels.Update) {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/calendarWithProgressHandler")
	defer tl.End()

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	img, err := i.service.GenerateCalendarPNG(ctx, userID, true)
	if err != nil {
		tl.AddError(err)

		if sErr := i.SendErrorMessage(ctx, chatID, err); sErr != nil {
			tl.AddError(sErr)
			return
		}

		return
	}

	caption := captionCalendarWithProgress
	name := fmt.Sprintf("calendar_with_progress_%s", time.Now().Format("2006-01-02"))

	if sErr := i.SendImagePNG(ctx, chatID, name, caption, img); sErr != nil {
		tl.AddError(sErr)
		return
	}
}

// errorsHandler - заглушка для bot.ErrorsHandler
func (i *implTelegramBot) errorsHandler(err error) {
	if strings.Contains(err.Error(), "getUpdates\": context canceled") {
		return
	}

	tl, _ := tracelog.Begin(context.Background(), "TGBOT/errorsHandler")
	tl.AddError(err)
	tl.End()
}

// debugHandler - заглушка для bot.DebugHandler
func (i *implTelegramBot) debugHandler(_ string, _ ...any) {
	tl, _ := tracelog.Begin(context.Background(), "TGBOT/debugHandler")
	tl.End()
}
