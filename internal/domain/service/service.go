package service

import (
	"context"
	"fmt"
	"trrader/internal/adapter"
	"trrader/internal/adapter/bybit"
	"trrader/internal/adapter/bybit/domain/linear"
	"trrader/internal/traidingview"
)

type service struct {
	ctx    context.Context
	Bybit  *adapter.Client
	Trview *traidingview.TraidingView
}

func New(tv *traidingview.TraidingView, cl *adapter.Client) *service {
	return &service{
		ctx:    context.Background(),
		Bybit:  cl,
		Trview: tv,
	}
}
func (s *service) StartTraiding() {

	for {
		alert := s.Trview.GetData()
		s.CreateOrder(alert)

	}

}
func (s *service) CreateOrder(a traidingview.Alert) error {
	cl := s.Bybit.Linear()
	p := &linear.PlaceActiveOrderParams{}
	p.SetParams(a.Ticker, bybit.Side(a.Side), 10.0)

	orderResponse, err := cl.Account().PlaceActiveOrder(s.ctx, p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderResponse)
	return nil
}
