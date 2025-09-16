package http

import (
	"ddd-timer-service/internal/pkg/stats_counter"
	"ddd-timer-service/models"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed src/index.html
var indexHTML []byte

func (i *implServerGin) rootHandler(c *gin.Context) {
	_, _ = c.Writer.Write(indexHTML)
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
}

func (i *implServerGin) statsHandler(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	fromDate, err := time.Parse(time.DateOnly, from)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid date format",
		})
		return
	}

	toDate, err := time.Parse(time.DateOnly, to)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid date format",
		})
		return
	}

	if !fromDate.Before(toDate) {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("fromDate must be before toDate"))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid date format",
		})
		return
	}

	user := &models.User{
		ServeFrom: fromDate,
		ServeTo:   toDate,
	}

	stats, err := stats_counter.NewStats(user, time.Now())
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(200, stats.PrettyShort())
}
