package main

import (
	"flag"

	"github.com/tencentad/martech/cmd/dumper/config"
	"github.com/tencentad/martech/cmd/dumper/targeting"
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

	dumper := targeting.NewDumper(config.Configuration.TargetingDumper)
	dumper.Dump()
}
