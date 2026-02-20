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

	logger.Info("1")
	logger.Warn("2")
	logger.Norm("3")
}
