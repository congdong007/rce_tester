package main

import (
	"rce_tester/config"
	"rce_tester/core"
)

func main() {
	cfg := config.ParseFlags()
	core.Run(cfg)
}
