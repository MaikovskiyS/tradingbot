package adapter

import "trrader/internal/adapter/bybit/linear"

type Client struct {
	linear linear.Client
}

func (c *Client) Linear() linear.Client {
	return c.linear
}

// NewRestClient - creates a new bybit rest client
// docs - https://bybit-exchange.github.io/docs/spot/#t-introduction
func NewRestClient(url, apiKey, apiSecret string) *Client {
	return &Client{
		linear: *linear.NewLinearClient(url, apiKey, apiSecret),
	}
}
