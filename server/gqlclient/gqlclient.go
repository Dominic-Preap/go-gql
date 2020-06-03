package gqlclient

import (
	"net/http"
	"reflect"

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

// Helper function to set variables into graphql request
//   - if using struct as a parameter when calling graphql, consider using it
// https://stackoverflow.com/questions/52562373/how-to-find-empty-struct-values-in-go-using-reflection
func setGraphQLVars(req *graphql.Request, x interface{}) {
	v := reflect.ValueOf(x)

	for i := 0; i < v.NumField(); i++ {
		// field := v.Type().Field(i).Name
		tag := v.Type().Field(i).Tag.Get("json")
		val := v.Field(i).Interface()

		// Check if the field is zero-valued, meaning it won't be updated
		if !reflect.DeepEqual(val, reflect.Zero(v.Field(i).Type()).Interface()) {
			// log.Printf("field --> %+v", tag)
			// log.Printf("value --> %+v", val)
			req.Var(tag, val)
		}
	}
}
