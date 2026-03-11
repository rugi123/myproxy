package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rugi123/myproxy/client/internal/common/config"
	"github.com/rugi123/myproxy/client/internal/common/logger"
	"github.com/rugi123/myproxy/client/internal/common/models"
)

type Server struct {
	config     *config.ServerConfig
	logger     *logger.Logger
	tunnelReq  chan []byte
	tunnelResp chan []byte
	tunnel     TunnelServer
	control    net.Conn
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

type TunnelServer struct {
	models.Tunnel
}

func New(cfg *config.ServerConfig, logger *logger.Logger) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		config:     cfg,
		logger:     logger,
		tunnelReq:  make(chan []byte),
		tunnelResp: make(chan []byte),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (s *Server) runTCP(handlers ...func(conn net.Conn)) error {
	port := fmt.Sprintf(":%d", s.config.TunnelPort)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("setup listener error: %v", err)
	}
	defer listener.Close()

	s.logger.Debug("TCP server starts on %s", port)

	//Graceful shutdown
	go func() {
		<-s.ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("accept error: %v", err)
			continue
		}

		s.wg.Add(1)
		go s.handleTCPConnect(conn, handlers...)
	}
}

func (s *Server) runHTTP(handler func(w http.ResponseWriter, r *http.Request)) error {
	port := fmt.Sprintf(":%d", s.config.BaseConfig.App.Port)
	server := &http.Server{
		Addr:         port,
		Handler:      http.HandlerFunc(handler),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//Graceful shutdown
	go func() {
		<-s.ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("HTTP server shutdown error: %v", err)
		}
	}()

	s.logger.Debug("HTTP server starts on %s", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}

func (s *Server) Start() error {
	errChan := make(chan error, 2)

	//запуск контроль сервера
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.runTCP(s.setupControl, s.controlHandler); err != nil {
			errChan <- err
		}
	}()

	//запуск тунеля
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.runTCP(s.tunnelHandler); err != nil {
			errChan <- err
		}
	}()

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		if err := s.runHTTP(s.apiHandler); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-s.ctx.Done():
		return nil
	}
}

func (s *Server) Shutdown() error {
	s.logger.Info("shutting down tunnel server...")
	s.cancel()
	s.wg.Wait()
	close(s.tunnelReq)
	close(s.tunnelResp)
	return nil
}
