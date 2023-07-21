package service

import (
	"context"
	"fmt"
	"strconv"
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
	par := &linear.GetSymbolInformationParams{
		Symbol: a.Ticker,
	}
	inforesp, err := cl.Market().GetSymbolInformation(s.ctx, par)
	price, _ := strconv.ParseFloat(inforesp.Result[0].LastPrice, 8)
	fmt.Println(price)
	//fmt.Println(inforesp.Result[0].LastPrice)
	p := &linear.PlaceActiveOrderParams{}
	p.SetParams(a.Ticker, bybit.Side(a.Side), 10.0, price)

	orderResponse, err := cl.Account().PlaceActiveOrder(s.ctx, p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderResponse)
	return nil
}
