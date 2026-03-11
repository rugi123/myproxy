package client

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/rugi123/myproxy/client/internal/common/config"
	"github.com/rugi123/myproxy/client/internal/common/logger"
	"github.com/rugi123/myproxy/client/internal/common/models"
)

type Client struct {
	serverIP   string
	serverPort int
	localPort  int
	control    net.Conn
	tunnel     ClientTunnel
	logger     *logger.Logger
}

func New(config *config.ClientConfig, logger *logger.Logger) *Client {
	return &Client{
		serverIP:   config.Server.IP,
		serverPort: config.Server.Port,
		localPort:  config.BaseConfig.App.Port,
		logger:     logger,
	}
}

func (c *Client) Connect() error {
	port := fmt.Sprintf("%d", c.serverPort)
	url := net.JoinHostPort(c.serverIP, port)
	control, err := net.Dial("tcp", url)
	if err != nil {
		return fmt.Errorf("connect server error: %v", err)
	}

	req := models.AuthRequest{
		Token: "123",
	}
	msg, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal auth req error: %v", err)
	}

	if _, err = control.Write(msg); err != nil {
		return fmt.Errorf("write auth req error: %v", err)
	}

	data := make([]byte, 1024)
	n, err := control.Read(data)
	if err != nil {
		return fmt.Errorf("read auth resp error: %v", err)
	}

	var resp models.AuthResponse
	if err := json.Unmarshal(data[:n], &resp); err != nil {
		return fmt.Errorf("marshal auth resp error: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("auth is not approved: %v", err)
	}

	c.control = control
	return nil
}
