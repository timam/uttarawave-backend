package server

import (
	"context"
	"github.com/spf13/viper"
	"net/http"
	"time"

	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"github.com/timam/uttaracloud-finance-backend/routers"
	"go.uber.org/zap"
)

var server *http.Server

func startServer() {
	router := routers.InitRouter()
	server = &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server: %v", zap.Error(err))
		}
	}()

	logger.Info("Server started on" + server.Addr)
}

func ReloadServer() {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal(
				"Failed to gracefully shutdown the server: %v",
				zap.Error(err))
		}
		logger.Info("Server exited")

	}

	startServer()
}

func StartServer() {
	startServer()
}
