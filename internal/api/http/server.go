package http

import (
	"context"
	"ddd-timer-service/internal/service"
	"ddd-timer-service/models"
	_ "embed"

	ginlogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server interface {
	Start(ctx context.Context, addr string) error
	Stop() error
}

type implServerGin struct {
	ok      bool
	e       *gin.Engine
	service *service.Service
}

func (i *implServerGin) Start(ctx context.Context, addr string) error {
	if i.ok == false {
		return models.ErrorNotInitialized
	}

	var err error

	go func() {
		err = i.e.Run(addr)
	}()

	<-ctx.Done()

	return err
}

func (i *implServerGin) Stop() error {
	return nil
}

func NewImplServerGin(service *service.Service) Server {
	gin.SetMode(gin.ReleaseMode)

	i := &implServerGin{
		ok:      true,
		service: service,
	}

	e := gin.New()

	e.Use(ginlogger.SetLogger(ginlogger.WithLogger(
		func(c *gin.Context, logger zerolog.Logger) zerolog.Logger {
			return service.Logger().With().Str("ip", c.ClientIP()).Str("user_agent", c.Request.UserAgent()).Logger()
		})))

	e.GET("/", i.rootHandler)
	e.NoRoute(i.rootHandler)

	e.GET("api/stats", i.statsHandler)

	i.e = e
	return i
}
