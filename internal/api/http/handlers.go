package http

import (
	"ddd-timer-service/internal/pkg/stats_counter"
	"ddd-timer-service/internal/pkg/tracelog"
	"ddd-timer-service/models"
	_ "embed"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

//go:embed src/index.html
var indexHTML []byte

func (i *implServerGin) rootHandler(c *gin.Context) {
	_, _ = c.Writer.Write(indexHTML)
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
}

func (i *implServerGin) statsHandler(c *gin.Context) {
	tl, _ := tracelog.Begin(c.Request.Context(), "HTTP/statsHandler")
	defer tl.End()

	from := c.Query("from")
	to := c.Query("to")

	fromDate, err := time.Parse(time.DateOnly, from)
	if err != nil {
		err = ErrInvalidDateFormat

		_ = c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{
			messageJsonKey: err.Error(),
		})

		tl.AddError(err, zerolog.WarnLevel)
		return
	}

	toDate, err := time.Parse(time.DateOnly, to)
	if err != nil {
		err = ErrInvalidDateFormat

		_ = c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{
			messageJsonKey: err.Error(),
		})

		tl.AddError(err, zerolog.WarnLevel)
		return
	}

	if !fromDate.Before(toDate) || fromDate.Round(time.Hour*24).Equal(toDate.Round(time.Hour*24)) {
		err = ErrBadDates

		_ = c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{
			messageJsonKey: err.Error(),
		})

		tl.AddError(err, zerolog.WarnLevel)
		return
	}

	user := &models.User{
		ServeFrom: fromDate,
		ServeTo:   toDate,
	}

	// Валидации для дат есть выше, поэтому ошибки здесь не будет
	stats, _ := stats_counter.NewStats(user, tl.StartedAt())

	c.JSON(200, stats.PrettyShort())

	tl.InfoWithDuration("stats generated", tracelog.String("date_from", fromDate.String()),
		tracelog.String("date_to", toDate.String()))
}
