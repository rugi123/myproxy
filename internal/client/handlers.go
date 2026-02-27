package client

/*

func (c *Client) tunnelHandler(conn net.Conn) {
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		c.logger.Info("read tunnel data error: %v", err)
		return
	}
	c.logger.Debug(string(data[:n]))



		port := fmt.Sprintf("%d", c.config.BaseConfig.App.Port)
		url := net.JoinHostPort("localhost", port)

		conn, err = net.Dial("tcp", url)
		if err != nil {
			c.logger.Info("setup connection error: %v", err)
			return
		}

		_, err = conn.Write(data[:n])
		if err != nil {
			c.logger.Info("write error: %v", err)
			return
		}



}

*/
