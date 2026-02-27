package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/rugi123/myproxy/client/internal/config"
	"github.com/rugi123/myproxy/client/internal/logger"
)

type TunnelServer struct {
	config     *config.ServerConfig
	logger     *logger.Logger
	tunnelChan chan []byte
}

func NewTunnelServer(cfg *config.ServerConfig, logger *logger.Logger) *TunnelServer {
	return &TunnelServer{
		config:     cfg,
		logger:     logger,
		tunnelChan: make(chan []byte),
	}
}

func runTcp(serverPort int, logger *logger.Logger, handlers ...func(conn net.Conn)) error {
	port := fmt.Sprintf(":%d", serverPort)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("setup listener error: %v", err)
	}
	defer listener.Close()

	logger.Debug("tunnel server starts on %s", port)

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("accept error: %v", err)
			continue
		}

		wg.Add(1)

		go func(conn net.Conn) {
			defer wg.Done()
			defer conn.Close()

			for _, handler := range handlers {
				handler(conn)
			}

		}(conn)
	}
}

func (s *TunnelServer) runHTTP(handler func(w http.ResponseWriter, r *http.Request)) error {
	port := fmt.Sprintf(":%d", s.config.BaseConfig.App.Port)

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}

func (s *TunnelServer) RunTunnel() error {
	err := runTcp(s.config.TunnelPort, s.logger, s.tunnelHandler)
	if err != nil {
		return err
	}
	return nil
}

func (s *TunnelServer) RunServer() error {
	err := s.runHTTP(s.serverHandler)
	if err != nil {
		return err
	}
	return nil
}
