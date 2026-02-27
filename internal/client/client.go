package client

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"

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

func (c *Client) RunTunnelClient() {
	for {
		// Пробуем подключиться и прочитать данные
		conn, err := net.Dial("tcp", "194.87.95.228:8080")
		if err != nil {
			continue
		}

		data := make([]byte, 1024)
		n, err := conn.Read(data)
		req := parseMessage(data[:n])
		fmt.Println(req)

	}
}

func parseMessage(message []byte) *http.Request {
	reader := bufio.NewReader(bytes.NewReader(message))
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}
	return req
}
