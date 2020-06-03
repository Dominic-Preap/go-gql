package gqlclient

import (
	"context"
	"errors"

	"github.com/machinebox/graphql"
)

// CountryData .
type CountryData struct {
	Country struct {
		Code    string `json:"code"`
		Name    string `json:"name"`
		Capital string `json:"capital"`
	} `json:"country"`
}

// GetCountry .
func (c *GQLClient) GetCountry(code string) (*CountryData, error) {
	q := /* GraphQL */ `
		query xxx($code: ID!) {
				country(code: $code) {
				code
				name
				capital
			}
		}	
	`
	// setGraphQLVars(req, cd)
	req := graphql.NewRequest(q)
	req.Var("code", code)

	var res = CountryData{}
	if err := c.Run(context.Background(), req, &res); err != nil {
		return nil, err
	}
	if (CountryData{} == res) {
		return nil, errors.New("Empty Object")
	}

	return &res, nil
}
