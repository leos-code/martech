package config_test

import (
	"os"
	"path"
	"testing"

	"github.com/tencentad/martech/cmd/web/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	p, _ := os.Getwd()
	p = path.Join(p, "..", "..", "..", "configs", "web.config.json")
	_, err := config.Load(p)
	assert.NoError(t, err)
}
