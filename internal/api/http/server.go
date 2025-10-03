package http

import (
	"context"
	"ddd-timer-service/internal/service"
	"ddd-timer-service/models"
	_ "embed"
	"io"

	"github.com/gin-gonic/gin"
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
	gin.DefaultWriter = io.Discard

	i := &implServerGin{
		ok:      true,
		service: service,
	}

	e := gin.New()

	e.Use(gin.Recovery())
	e.Use(accessLogMW())

	e.GET("/", i.rootHandler)
	e.NoRoute(i.rootHandler)

	e.GET("api/stats", i.statsHandler)

	i.e = e
	return i
}
