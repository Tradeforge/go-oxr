package client

import (
	"log/slog"

	"go.tradeforge.dev/oxr/internal/client"
)

const (
	OpenExchangeRatesAPIURL = "https://openexchangerates.org/api"
	authSchemeToken         = "Token"
)

type Client struct {
	client.Client

	CurrencyClient
}

func NewClient(
	apiURL string,
	apiKey string,
	options ...Option,
) *Client {
	c := client.New(apiURL, slog.Default())
	c.SetAuthScheme(authSchemeToken)
	c.SetAuthToken(apiKey)
	for _, option := range options {
		option(c)
	}
	return &Client{
		Client:         *c,
		CurrencyClient: CurrencyClient{Client: c},
	}
}
