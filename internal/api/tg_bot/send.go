package tg_bot

import (
	"bytes"
	"context"
	"ddd-timer-service/internal/pkg/tracelog"
	"fmt"

	"github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
)

func (i *implTelegramBot) SendMessage(ctx context.Context, chatID int64, pm botmodels.ParseMode, message string) error {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/SendMessage")
	defer tl.End()

	if _, err := i.pkgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      message,
		ParseMode: pm,
	}); err != nil {
		tl.AddError(err)
		return err
	}

	return nil
}

func (i *implTelegramBot) SendErrorMessage(ctx context.Context, chatID int64, sErr error, message ...string) error {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/SendErrorMessage")
	defer tl.End()

	var text string
	if len(message) == 0 {
		text = sErr.Error()
	} else {
		text = fmt.Sprintf("%s: %s", message[0], sErr)
	}

	if _, err := i.pkgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		tl.AddError(err)
		return err
	}

	return nil
}

func (i *implTelegramBot) SendImagePNG(ctx context.Context, chatID int64, name, caption string, img []byte) error {
	tl, ctx := tracelog.Begin(ctx, "TGBOT/SendImagePNG")
	defer tl.End()

	media := &botmodels.InputMediaPhoto{
		Media:           fmt.Sprintf("attach://%s.png", name),
		Caption:         caption,
		MediaAttachment: bytes.NewReader(img),
	}

	if _, err := i.pkgBot.SendMediaGroup(ctx, &bot.SendMediaGroupParams{
		ChatID: chatID,
		Media: []botmodels.InputMedia{
			media,
		},
	}); err != nil {
		tl.AddError(err)
		return err
	}

	return nil
}
