package calendar_drawer

import (
	_ "embed"
	"image"
	"time"
)

type CalendarDrawer struct{}

func NewCalendarDrawer() *CalendarDrawer {
	return &CalendarDrawer{}
}

func (c *CalendarDrawer) BySeasonsPNG(from, to time.Time) ([]byte, image.Image, error) {
	seasons := NewCalendar(from, to).Seasons()

	return seasons.PNG()
}

func (c *CalendarDrawer) BySeasonsWithProgressPNG(from, to, now time.Time) ([]byte, image.Image, error) {
	seasons := NewCalendar(from, to).Seasons()

	return seasons.PNGWithProgressMask(from, to, now)
}
