package client

import (
	"io"
	"log"
	"net"
	"time"

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

func (c *Client) run() error {
	for {
		// Пробуем подключиться и прочитать данные
		conn, err := net.Dial("tcp", "194.87.95.228:8080")
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		// Читаем все что придет и выводим
		io.Copy(log.Writer(), conn)
		conn.Close()

		time.Sleep(5 * time.Second)
	}
}

func (c *Client) RunTunnelClient() error {
	if err := c.run(); err != nil {
		return err
	}
	return nil
}
