package server

import "net"

func (s *TunnelServer) tunnelHandler(conn net.Conn) {
	buffer := make([]byte, 1024)
	buffer = <-s.tunnelChan

	_, err := conn.Write(buffer)
	if err != nil {
		s.logger.Error("write message error: %v", err)
	}
}

func (s *TunnelServer) serverHandler(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		s.logger.Error("read message error: %v", err)
	}

	s.tunnelChan <- buffer[:n]
}
