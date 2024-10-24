package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/timam/uttarawave-backend/api/routers"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	ctx        context.Context
	cancel     context.CancelFunc
}

func InitializeServer() (*Server, error) {
	router := routers.InitRouter()
	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		httpServer: &http.Server{
			Addr:    ":" + viper.GetString("server.port"),
			Handler: router,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	logger.Info("Server initialized successfully")
	return server, nil
}

func (s *Server) RunServer() error {
	logger.Info("Server starting on " + s.httpServer.Addr)

	go func() {
		<-s.ctx.Done()
		if err := s.GracefulShutdown(); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	newServer, err := InitializeServer()
	if err != nil {
		return fmt.Errorf("failed to initialize new server: %w", err)
	}

	if err := newServer.RunServer(); err != nil {
		return fmt.Errorf("failed to run new server: %w", err)
	}

	*s = *newServer
	return nil
}

func (s *Server) ShutdownServer() error {
	s.cancel()
	err := s.GracefulShutdown()
	logger.SyncLogger()
	return err
}
