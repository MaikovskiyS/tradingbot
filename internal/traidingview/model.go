package traidingview

import (
	"errors"
	"fmt"
	"time"
)

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
type AdxStrategyAlert struct {
	Adx        SideAlert
	Volatility SideAlert
	Volume     ActionAlert
}
type BreakOutAlert struct {
	MainTrend  SideAlert
	Trend      SideAlert
	Volatility ActionAlert
	Volume     ActionAlert
}

func (as *AdxStrategyAlert) Validate() error {
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
func (as *BreakOutAlert) Validate() error {
	timestamp := time.Now().Unix()
	fmt.Println("time now:", timestamp)
	if as.Trend.Side == "" && as.Trend.Ticker == "" {
		return errors.New("trend cant be empty")
	}
	if as.Trend.TimeStamp < timestamp-60 {
		fmt.Println("trend time:", as.Trend.TimeStamp)
		fmt.Println("real time:", timestamp)
		return errors.New("trend timestamp invalid")
	}
	if as.MainTrend.Side == "" && as.MainTrend.Ticker == "" {
		return errors.New("maintrend cant be empty")
	}
	if as.MainTrend.TimeStamp < timestamp-60 {
		fmt.Println("maintrend time:", as.MainTrend.TimeStamp)
		fmt.Println("real time:", timestamp)
		return errors.New("maintrend timestamp invalid")
	}
	if as.Trend.Side != as.MainTrend.Side {
		return errors.New("indefinite trend")
	}
	if as.Volatility.Action == "" && as.Volatility.Ticker == "" {
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
	if as.Volume.TimeStamp < timestamp-60 {
		fmt.Println("volume time:", as.Volume.TimeStamp)
		fmt.Println("real time:", timestamp)
		return errors.New("volume timestamp invalid")
	}
	return nil
}
