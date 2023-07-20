package service

import (
	"context"
	"fmt"

	"trrader/internal/adapter"
	"trrader/internal/adapter/bybit/domain/linear"
	"trrader/internal/traidingview"
)

type service struct {
	Bybit  *adapter.Client
	Trview *traidingview.TraidingView
}

func New(tv *traidingview.TraidingView, cl *adapter.Client) *service {
	return &service{
		Bybit:  cl,
		Trview: tv,
	}
}
func (s *service) StartTraiding() {
	bb := s.Bybit.Linear()

	ctx := context.Background()
	p := &linear.OrderBookParams{}
	for {
		alert := s.Trview.GetData()
		p.Symbol = alert.Ticker
		book, err := bb.Market().GetOrderBook(ctx, p)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(book)
	}

}
