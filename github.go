package main

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type authedTransport struct {
	wrapped http.RoundTripper
	token   string
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.token)
	return t.wrapped.RoundTrip(req)
}

func client(token string) graphql.Client {
	return graphql.NewClient("https://api.github.com/graphql",
		&http.Client{Transport: &authedTransport{
			wrapped: http.DefaultTransport,
			token:   token,
		}})
}

//go:generate go run github.com/Khan/genqlient
