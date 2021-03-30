package main

import (
	"flag"

	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/tencentad/martech/cmd/web/router"
	log "github.com/sirupsen/logrus"
)

var (
	configPath = flag.String("config_path", "web.config.json", "")
)

func main() {
	flag.Parse()

	configuration, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := config.Setup(configuration); err != nil {
		log.Fatal(err)
	}

	r, err := router.SetupRouter(configuration.Web)
	if err != nil {
		log.Fatal(err)
	}

	if err = handler.SetupHandler(configuration); err != nil {
		log.Fatal(err)
	}

	if err := r.Run(configuration.Web.ServerAddress); err != nil {
		log.Fatal(err)
	}
}
