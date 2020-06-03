package gqlclient

import (
	"net/http"

	"github.com/machinebox/graphql"
	"github.com/my/app/server/config"
)

type transport struct {
	headers map[string]string
	base    http.RoundTripper
}

// GQLClient .
type GQLClient struct {
	*graphql.Client
}

// Init .
func Init(env *config.EnvConfig) *GQLClient {
	client := graphql.NewClient("https://countries.trevorblades.com", graphql.WithHTTPClient(&http.Client{
		Transport: &transport{
			headers: map[string]string{
				"Cache-Control": "no-cache",
				"Authorization": "Bearer ...",
			},
		},
	}))
	return &GQLClient{Client: client}
}

// RoundTrip handle each request with the same headers
// https://stackoverflow.com/questions/54088660/add-headers-for-each-http-request-using-client
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}
