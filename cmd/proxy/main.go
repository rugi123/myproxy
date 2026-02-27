package main

import (
	"fmt"
	"os"

	"github.com/rugi123/myproxy/client/internal/client"
	"github.com/rugi123/myproxy/client/internal/config"
	"github.com/rugi123/myproxy/client/internal/logger"
)

func main() {
	cfg, err := config.LoadClient("./internal/config/")
	if err != nil {
		fmt.Printf("load conf error: %v", err)
		os.Exit(1)
	}

	log := logger.New(logger.Level(cfg.BaseConfig.LogLevel), os.Stdout, make(chan logger.Entry))

	go log.Run()

	server := client.NewClient(cfg, log)

	if err := server.RunTunnelClient(); err != nil {
		log.Fatal("run server error: %v", err)
	}
}
