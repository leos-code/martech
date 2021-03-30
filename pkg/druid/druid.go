package druid

import (
	"github.com/awatercolorpen/godruid"
	"net/http"
	"time"
)

func GetDruidClient() *godruid.Client {
	return &godruid.Client{
		Url: "http://9.134.188.241:7888",
		Debug: true,
		HttpClient: &http.Client{Timeout: 60 * time.Second},
	}
}
