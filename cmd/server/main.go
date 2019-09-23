package main

import (
	"github.com/rejlersembriq/hooked/pkg/repository/memory"
	"github.com/rejlersembriq/hooked/pkg/router"
	"github.com/rejlersembriq/hooked/pkg/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"time"
)

func main() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig = config

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      server.New(router.New(), memory.New()),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	zap.L().Info("Starting server", zap.String("address", srv.Addr))
	if err := srv.ListenAndServe(); err != nil {
		zap.L().Info("Error serving http.", zap.String("error", err.Error()))
	}
}
