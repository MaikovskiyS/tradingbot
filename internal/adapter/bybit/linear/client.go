package linear

import (
	"trrader/internal/adapter/bybit/domain/linear"
	"trrader/internal/adapter/bybit/linear/account"
	"trrader/internal/adapter/bybit/linear/market"
)

type Client struct {
	market  linear.MarketInterface
	account linear.AccountInterface
}

func (c *Client) Market() linear.MarketInterface {
	return c.market
}

func (c *Client) Account() linear.AccountInterface {
	return c.account
}

func NewLinearClient(url, apiKey, apiSecret string) *Client {
	return &Client{
		market:  market.NewLinearMarketClient(url, apiKey, apiSecret),
		account: account.NewLinearAccountClient(url, apiKey, apiSecret),
	}
}
