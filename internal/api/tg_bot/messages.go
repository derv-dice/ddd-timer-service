package tg_bot

import (
	"fmt"

	"github.com/go-telegram/bot"
)

const (
	mustRegisterMessage = "Для работы бота требуется ввести даты начала/окончания службы в формате `01.01.2024 01.01.2025`, где первая дата это начало службы, а вторая окончания"
	helpMessage         = "Чтобы изменить даты начала/окончания службы, отправь сообщение в формате `01.01.2024 01.01.2025`, где первая дата это начало службы, а вторая окончания"
	basicUsageMessage   = "Чтобы узнать статистику используй команду /stats"

	captionCells                = ""
	captionCalendar             = `Посезонный календарь на весь период службы. Можешь распечатать его и отмечать дни вручную, либо вызвать команду /calendar_with_progress и получить картинку уже заполненного календаря`
	captionCalendarWithProgress = `Посезонный календарь на весь период службы с отметками о прошедших днях. Чтобы получить такой же, но без отметок, вызови команду /calendar`
)

func newHelloMessage(firstname string) string {
	return fmt.Sprintf("Привет, *%s*", bot.EscapeMarkdown(firstname))
}
