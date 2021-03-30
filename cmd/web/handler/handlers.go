package handler

import "github.com/tencentad/martech/cmd/web/config"

func SetupHandler(config *config.Configuration) error {
	materialFileOption = config.MaterialFile
	return nil
}
