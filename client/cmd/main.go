package main

import (
	"fmt"

	"github.com/rugi123/myproxy/client/internal/config"
	"github.com/rugi123/myproxy/client/internal/logger"
)

func main() {
	err, cfg := config.Load("internal/config/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)

	var level logger.Level = logger.Level(cfg.LogLevel)
	log := logger.New(level)

	log.Debug("a")
	log.Info("b")
	log.Warn("c")
	log.Error("d")
	log.Fatal("d")
}
