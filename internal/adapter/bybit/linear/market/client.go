package market

import (
	"trrader/internal/adapter/bybit/transport"
	"trrader/internal/adapter/bybit/transport/http"
)

type LinearMarketClient struct {
	Transporter transport.Transporter
}

func NewLinearMarketClient(url, apiKey, apiSecret string) *LinearMarketClient {
	transporter := http.New(url, apiKey, apiSecret)
	return &LinearMarketClient{
		Transporter: transporter,
	}
}
