package calendar_drawer

import (
	"context"
	"ddd-timer-service/internal/pkg/img_cache"
	"ddd-timer-service/internal/pkg/tracelog"
	_ "embed"
	"fmt"
	"time"
)

type CalendarDrawer struct {
	maxYears int

	cache *img_cache.Cache
}

func NewCalendarDrawer(ctx context.Context, maxYears, maxCacheSizeMB int) *CalendarDrawer {
	cacheSize := maxCacheSizeMB * 1024 * 1024
	cache := img_cache.NewCache(cacheSize)

	c := &CalendarDrawer{
		maxYears: maxYears,
		cache:    cache,
	}

	go c.backgroundCacheCleaner(ctx)

	return c
}

func (c *CalendarDrawer) BySeasonsPNG(from, to time.Time, disableLimits bool) ([]byte, error) {
	seasons := NewCalendar(from, to).Seasons()

	if !disableLimits {
		if err := c.validateDates(from, to); err != nil {
			return nil, err
		}
	}

	key := c.newUniqueKeyFromDates(from, to)

	cachedImgBytes := c.cache.Get(key, imgTCalendarBySeasonsPNG)
	if cachedImgBytes != nil {
		return cachedImgBytes, nil
	}

	imgBytes, _, err := seasons.PNG()
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, imgTCalendarBySeasonsPNG, imgBytes)

	return imgBytes, nil
}

func (c *CalendarDrawer) BySeasonsWithProgressPNG(from, to, now time.Time, disableLimits bool) ([]byte, error) {
	seasons := NewCalendar(from, to).Seasons()

	if !disableLimits {
		if err := c.validateDates(from, to); err != nil {
			return nil, err
		}
	}

	key := c.newUniqueKeyFromDates(from, to)

	cachedImgBytes := c.cache.Get(key, imgTCalendarBySeasonsWithProgressPNG)
	if cachedImgBytes != nil {
		return cachedImgBytes, nil
	}

	imgBytes, _, err := seasons.PNGWithProgressMask(from, to, now)
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, imgTCalendarBySeasonsWithProgressPNG, imgBytes)

	return imgBytes, nil
}

// backgroundCacheCleaner - очистка кеша раз в сутки в 00:00
func (c *CalendarDrawer) backgroundCacheCleaner(ctx context.Context) {
	tl, ctx := tracelog.Begin(ctx, "calendarDrawer/backgroundCacheCleaner")
	defer tl.End()

	for {
		now := time.Now()

		next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		tl.Info(fmt.Sprintf("next cleaning starts at %s", next.Format(time.DateTime)))

		timer := time.NewTimer(next.Sub(now))

		select {
		case <-timer.C:
			c.cache.Clear()
			tl.Info("cache is cleaned")

		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
}

func (c *CalendarDrawer) validateDates(from, to time.Time) error {
	if c.maxYears != 0 && int(to.Sub(from).Hours()/24/365) > c.maxYears {
		return fmt.Errorf("временной период слишком большой. maxYearsPeriod=%d", c.maxYears)
	}

	return nil
}

func (c *CalendarDrawer) newUniqueKeyFromDates(from, to time.Time) string {
	return from.Round(time.Hour*24).Format(keyTimeFormat) + to.Round(time.Hour*24).Format(keyTimeFormat)
}
