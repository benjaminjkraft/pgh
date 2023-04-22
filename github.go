package main

import (
	"net/http"
	"os"

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

func mustGetToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("need to set GITHUB_TOKEN")
	}
	return token
}

//go:generate go run github.com/Khan/genqlient
