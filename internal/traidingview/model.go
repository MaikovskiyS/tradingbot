package traidingview

import (
	"errors"
	"fmt"
	"time"
)

type StrategyAlert struct {
	Adx        SideAlert
	Volatility SideAlert
	Volume     ActionAlert
}
type SideAlert struct {
	Ticker    string
	Side      string
	TimeStamp int64
}
type ActionAlert struct {
	Ticker    string
	Action    string
	TimeStamp int64
}

func (as *StrategyAlert) Validate() error {
	timestamp := time.Now().Unix()
	fmt.Println("time now:", timestamp)
	if as.Adx.Side == "" && as.Adx.Ticker == "" {
		return errors.New("adx cant be empty")
	}
	if as.Adx.TimeStamp < timestamp-60 {
		fmt.Println("adx time:", as.Adx.TimeStamp)
		fmt.Println("real time:", timestamp)
		return errors.New("adx timestamp invalid")
	}
	if as.Volatility.Side == "" && as.Volatility.Ticker == "" {
		return errors.New("volatility cant be empty")
	}
	if as.Volatility.TimeStamp < timestamp-60 {
		fmt.Println("volatility time:", as.Volatility.TimeStamp)
		fmt.Println("real time:", timestamp)
		return errors.New("volatility timestamp invalid")
	}
	if as.Volume.Action == "" && as.Volume.Ticker == "" {
		return errors.New("volume cant be empty")
	}
	return nil
}
