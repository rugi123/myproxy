package main

import (
	"fmt"

	"github.com/rugi123/myproxy/client/internal/config"
)

func main() {
	err, cfg := config.Load("internal/config/")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}
