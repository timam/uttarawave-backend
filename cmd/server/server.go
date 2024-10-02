package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"github.com/timam/uttaracloud-finance-backend/routers"
	"go.uber.org/zap"
)

var (
	server       *http.Server
	serverCtx    context.Context
	serverCancel context.CancelFunc
	mu           sync.Mutex
)

func initServerContext() {
	mu.Lock()
	defer mu.Unlock()
	serverCtx, serverCancel = context.WithCancel(context.Background())
}

func startServer() *http.Server {
	router := routers.InitRouter()
	return &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: router,
	}
}

func StartServer() {
	initServerContext()
	server = startServer()
	logger.Info("Server started on " + server.Addr)

	go func() {
		<-serverCtx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", zap.Error(err))
	}
}

func ReloadServer() {
	mu.Lock()
	defer mu.Unlock()

	// Cancel the existing server context
	if serverCancel != nil {
		serverCancel()
	}

	// Create a new context for the new server instance
	serverCtx, serverCancel = context.WithCancel(context.Background())

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Failed to gracefully shutdown the server", zap.Error(err))
		}
		logger.Info("Server exited")
	}

	// Start the new server instance
	go StartServer()
}

func ShutdownServer() error {
	mu.Lock()
	defer mu.Unlock()

	if server != nil && serverCancel != nil {
		serverCancel() // Cancel the server context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
	return nil
}
