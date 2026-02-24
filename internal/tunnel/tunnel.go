package tunnel

import (
	"bufio"
	"net"

	"github.com/rugi123/myproxy/client/internal/logger"
)

func HandleConnect(conn net.Conn, logger *logger.Logger) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		logger.Info("got message: %s", message)
	}
}
