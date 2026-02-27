package client

import (
	"fmt"
	"net"
	"sync"

	"github.com/rugi123/myproxy/client/internal/config"
	"github.com/rugi123/myproxy/client/internal/logger"
)

type Client struct {
	config *config.ClientConfig
	logger *logger.Logger
}

func NewClient(config *config.ClientConfig, logger *logger.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

func (c *Client) run(handlers ...func(net.Conn)) error {
	port := fmt.Sprintf("%d", c.config.Server.Port)
	url := net.JoinHostPort(c.config.Server.IP, port)

	conn, err := net.Dial("tcp", url)
	if err != nil {
		return fmt.Errorf("setup connection error: %v", err)
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)

	go func(conn net.Conn) {
		for _, handler := range handlers {
			handler(conn)
		}
	}(conn)

	return nil
}

func (c *Client) RunTunnelClient() error {
	if err := c.run(c.tunnelHandler); err != nil {
		return err
	}
	return nil
}
