package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/routers"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	ctx        context.Context
	cancel     context.CancelFunc
}

func InitServer() *Server {
	router := routers.InitRouter()
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + viper.GetString("server.port"),
			Handler: router,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) RunServer() error {
	logger.Info("Server starting on " + s.httpServer.Addr)

	go func() {
		<-s.ctx.Done()
		if err := s.GracefulShutdown(); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *Server) GracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) ReloadServer() error {
	if err := s.GracefulShutdown(); err != nil {
		return fmt.Errorf("failed to gracefully shutdown the server: %w", err)
	}
	logger.Info("Server exited")

	newServer := InitServer()
	go newServer.RunServer()

	*s = *newServer
	return nil
}

func (s *Server) ShutdownServer() error {
	s.cancel()
	return s.GracefulShutdown()
}
