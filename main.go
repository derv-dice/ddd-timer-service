package main

import (
	"context"
	"ddd-timer-service/config"
	"ddd-timer-service/internal/app"
	"ddd-timer-service/internal/pkg/tracelog"
	"ddd-timer-service/internal/pkg/tracelog/exporters"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	logger.Info().Msg("reading config")
	conf, err := config.ReadConfig(config.DefaultConfigFilename)
	if err != nil {
		logger.Fatal().Err(err).Msg("reading config failed")
	}

	level, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		logger.Fatal().Err(err).Msg("parsing log_level failed")
	}
	
	zerolog.SetGlobalLevel(level)

	tp, _ := exporters.ExporterDevNull()
	tracelog.SetLogger(logger)
	otel.SetTracerProvider(tp)

	wg := new(sync.WaitGroup)

	logger.Info().Msg("running service")
	err = app.Run(ctx, wg, *conf, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("run application failed")
	}

	logger.Info().Msg("service is running")

	logger.Info().Msg(app.StartupMessage)
	logger.Info().Msgf("http server is running on http://%s", conf.Http.Addr)

	wg.Wait()

	logger.Info().Msg("service stopped")
}
