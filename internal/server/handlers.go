package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/rugi123/myproxy/client/internal/common/models"
)

func (s *Server) handleTCPConnect(conn net.Conn, handlers ...func(conn net.Conn)) {
	for _, handler := range handlers {
		handler(conn)
	}
}

func (s *Server) apiHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := httputil.DumpRequest(r, true)
	if err != nil {
		s.logger.Info("read req error: %v", err)
	}
	s.tunnelReq <- msg

	msg = <-s.tunnelResp
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(msg))
}

func (s *Server) tunnelHandler(conn net.Conn) {
	resp := make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		s.logger.Info("read tunnel error: %v", err)
		return
	}
	s.tunnelResp <- resp[:n]

	req := make([]byte, 1024)
	req = <-s.tunnelReq

	if _, err := conn.Write(req); err != nil {
		s.logger.Info("write tunnel error: %v", err)
		return
	}
}

func (s *Server) controlHandler(conn net.Conn) {

}

func (s *Server) setupControl(conn net.Conn) {
	if s.control != nil {
		return
	}
	msg := make([]byte, 1024)
	n, err := conn.Read(msg)
	if err != nil {
		s.logger.Info("read control req error: %v", err)
		return
	}

	var req models.AuthRequest
	if err := json.Unmarshal(msg[:n], &req); err != nil {
		s.logger.Info("unmarshal control req error: %v", err)
		return
	}

	res := models.AuthResponse{
		Success: true,
		Message: "bad token",
	}
	if req.Token != "123" {
		s.logger.Info("bad token")
		res.Success = false
	} else {
		s.control = conn
	}

	msg, err = json.Marshal(res)
	if err != nil {
		s.logger.Info("marshal control resp error: %v", err)
		return
	}

	if _, err = conn.Write(msg); err != nil {
		s.logger.Info("write control resp error: %v", err)
		return
	}
}
