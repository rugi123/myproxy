package baseServer

import (
	"fmt"
	"net"

	"github.com/rugi123/myproxy/client/internal/config"
	"github.com/rugi123/myproxy/client/internal/logger"
)

type Handler func(conn net.Conn, logger *logger.Logger)

type Server struct {
	config  *config.AppConfig
	logger  *logger.Logger
	handler Handler
}

func NewServer(cfg *config.AppConfig, logger *logger.Logger, handler Handler) *Server {
	return &Server{
		config:  cfg,
		logger:  logger,
		handler: handler,
	}
}

func (s *Server) RunServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("setup listener error: %v", err)
	}

	defer listener.Close()

	s.logger.Debug("server starts on %d", s.config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Info("")
			continue
		}

		s.handler(conn, s.logger)
	}
}
