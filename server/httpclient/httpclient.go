package httpclient

import (
	"github.com/go-resty/resty/v2"
	"github.com/my/app/server/config"
)

// HTTPClient .
type HTTPClient struct {
	*resty.Client
}

// Init .
func Init(env *config.EnvConfig) *HTTPClient {
	// Create a Resty Client
	client := resty.New().
		SetHostURL("https://reqres.in/api/").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer ...")

	return &HTTPClient{Client: client}
}
