package calendar_drawer

import (
	"context"
	"ddd-timer-service/internal/pkg/img_cache"
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

	go c.runCacheCleaner(ctx)

	return c
}

func (c *CalendarDrawer) BySeasonsPNG(userID int64, from, to time.Time, disableLimits bool) ([]byte, error) {
	seasons := NewCalendar(from, to).Seasons()

	if !disableLimits {
		if err := c.validateDates(from, to); err != nil {
			return nil, err
		}
	}

	cachedImgBytes := c.cache.Get(userID, imgTCalendarBySeasonsPNG)
	if cachedImgBytes != nil {
		return cachedImgBytes, nil
	}

	imgBytes, _, err := seasons.PNG()
	if err != nil {
		return nil, err
	}

	c.cache.Set(userID, imgTCalendarBySeasonsPNG, imgBytes)

	return imgBytes, nil
}

func (c *CalendarDrawer) BySeasonsWithProgressPNG(userID int64, from, to, now time.Time, disableLimits bool) ([]byte, error) {
	seasons := NewCalendar(from, to).Seasons()

	if !disableLimits {
		if err := c.validateDates(from, to); err != nil {
			return nil, err
		}
	}

	cachedImgBytes := c.cache.Get(userID, imgTCalendarBySeasonsWithProgressPNG)
	if cachedImgBytes != nil {
		return cachedImgBytes, nil
	}

	imgBytes, _, err := seasons.PNGWithProgressMask(from, to, now)
	if err != nil {
		return nil, err
	}

	c.cache.Set(userID, imgTCalendarBySeasonsWithProgressPNG, imgBytes)

	return imgBytes, nil
}

// runCacheCleaner - очистка кеша раз в сутки в 00:00
func (c *CalendarDrawer) runCacheCleaner(ctx context.Context) {
	for {
		now := time.Now()

		next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		timer := time.NewTimer(next.Sub(now))

		select {
		case <-timer.C:
			c.cache.Clear()

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
