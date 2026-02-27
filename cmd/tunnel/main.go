package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/rugi123/myproxy/client/internal/config"
	"github.com/rugi123/myproxy/client/internal/logger"
	"github.com/rugi123/myproxy/client/internal/server"
)

func main() {
	cfg, err := config.LoadServer("./internal/config/")
	if err != nil {
		fmt.Printf("load conf error: %v", err)
		os.Exit(1)
	}

	log := logger.New(logger.Level(cfg.BaseConfig.LogLevel), os.Stdout, make(chan logger.Entry))

	go log.Run()

	log.Info("test")

	server := server.NewTunnelServer(cfg, log)

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(2)

	go server.RunTunnel()

	go server.RunServer()

}
