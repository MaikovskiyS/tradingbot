package service

import (
	"context"
	"fmt"
	"sort"
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
	params, err := s.SetActiveOrderParams(a)
	orderResponse, err := cl.Account().PlaceActiveOrder(s.ctx, params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderResponse)
	return nil
}
func (s *service) SetActiveOrderParams(a traidingview.Alert) (p *linear.PlaceActiveOrderParams, err error) {
	cl := s.Bybit.Linear()
	//price
	inforesp, err := cl.Market().GetSymbolInformation(s.ctx, &linear.GetSymbolInformationParams{Symbol: a.Ticker})
	if err != nil {
		return &linear.PlaceActiveOrderParams{}, err
	}
	price, err := strconv.ParseFloat(inforesp.Result[0].LastPrice, 8)
	if err != nil {
		return &linear.PlaceActiveOrderParams{}, err
	}
	stoploss, err := s.GetStopLoss(a)
	if err != nil {
		return &linear.PlaceActiveOrderParams{}, err
	}

	p.SetParams(a.Ticker, bybit.Side(a.Side), 10.0, stoploss, price)
	return
}
func (s *service) GetStopLoss(a traidingview.Alert) (float64, error) {
	cl := s.Bybit.Linear()
	//klines for stoploss
	resp, err := cl.Market().QueryKline(s.ctx, &linear.QueryKlineParams{Symbol: a.Ticker, Interval: bybit.Interval5Min, From: 1500000, Limit: 5})
	if err != nil {
		return 0.0, err
	}
	//TODO STOPLOSS.....
	var b linear.LowByPrice
	klines := resp.Result
	b = klines
	sort.Sort(b)
	stop := b[0].Low
	return stop, nil

}
