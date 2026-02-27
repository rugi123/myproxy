package main

import (
	"fmt"
	"os"

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

	server := server.NewTunnelServer(cfg, log)

	if err := server.RunTunnel(); err != nil {
		log.Fatal("run tunnel server error: %v", err)
	}

	if err := server.RunServer(); err != nil {
		log.Fatal("run server error: %v", err)
	}
}
