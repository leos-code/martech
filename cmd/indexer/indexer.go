package main

import (
	"flag"

	"github.com/tencentad/martech/cmd/indexer/config"
	"github.com/tencentad/martech/cmd/indexer/process"
	log "github.com/sirupsen/logrus"
)

var (
	configPath = flag.String("config_path", "", "")
)

func main() {
	flag.Parse()
	err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config, err: %v", err)
	}

	indexer := process.NewIndexProcessor(config.Configuration.Index)
	indexer.Process()
}
