package account

import (
	"trrader/internal/adapter/bybit/transport"
	"trrader/internal/adapter/bybit/transport/http"
)

type LinearAccountClient struct {
	Transporter transport.Transporter
}

func NewLinearAccountClient(url, apiKey, apiSecret string) *LinearAccountClient {
	transporter := http.New(url, apiKey, apiSecret)
	return &LinearAccountClient{
		Transporter: transporter,
	}
}
