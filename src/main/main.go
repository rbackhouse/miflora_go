package main

import (
	"os"

	logger "github.com/sirupsen/logrus"

	"potpie.org/miflora/src/config"
	"potpie.org/miflora/src/scanner"
)

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		logger.Warn("No command provided - specifiy \"scan\" or \"discover\"")
		return
	}
	if argsWithoutProg[0] == "scan" {
		cfg := config.NewConfig()
		scanner.Scan(cfg)
	} else if argsWithoutProg[0] == "discover" {
		scanner.Discover()
	} else {
		logger.Warn("Invalid command")
	}
}
