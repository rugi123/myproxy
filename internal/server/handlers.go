package server

import (
	"net"
	"net/http"
	"net/http/httputil"
)

func (s *TunnelServer) tunnelHandler(conn net.Conn) {
	buffer := make([]byte, 1024)
	buffer = <-s.tunnelChan

	_, err := conn.Write(buffer)
	if err != nil {
		s.logger.Error("write message error: %v", err)
	}
}

/*
func (s *TunnelServer) serverHandler(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		s.logger.Error("read message error: %v", err)
	}

	s.tunnelChan <- buffer[:n]
}
*/

func (s *TunnelServer) serverHandler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true) // true = включая тело
	if err != nil {
		http.Error(w, "Error dumping request", http.StatusInternalServerError)
		return
	}
	s.logger.Debug(string(dump))
	s.tunnelChan <- dump
}
