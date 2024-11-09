package main

import (
	"AuthAPI/cfg"
	"AuthAPI/internal/infra"
	
	"log"
	"os"
)

func main() {
	cfg, err := cfg.LoadCfg()
	if err != nil {
		log.Fatal("Failed to load environment configuration.")
		os.Exit(1)
	}

	infra.Run(cfg)

}