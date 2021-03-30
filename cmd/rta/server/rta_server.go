package main

import (
	"flag"
	"net/http"

	"github.com/tencentad/martech/cmd/rta/server/config"
	"github.com/tencentad/martech/cmd/rta/server/logic"
	"github.com/tencentad/martech/pkg/common/loader"
	"github.com/tencentad/martech/pkg/common/metricutil"
	"github.com/tencentad/martech/pkg/retrieval"
	log "github.com/sirupsen/logrus"
)

var (
	configFile     = flag.String("config_file", "", "")
	serverAddress  = flag.String("server_address", ":80", "")
	metricsAddress = flag.String("metrics_address", ":8080", "")
)

func serveHttp(serverAddress string) error {
	service := logic.NewRTAService(config.Configuration.RTAService)
	http.Handle("/rta", service)
	return http.ListenAndServe(serverAddress, nil)
}

func main() {
	flag.Parse()

	if err := config.Load(*configFile); err != nil {
		log.Fatal(err)
	}
	loader.StartDoubleBufferLoad(1)

	if config.Configuration.Retrieval != nil {
		if err := retrieval.GetServiceImpl().Start(); err != nil {
			log.Fatal(err)
		}
	}

	_ = metricutil.ServeMetrics(*metricsAddress)
	if err := serveHttp(*serverAddress); err != nil {
		log.Fatal(err)
	}
}
