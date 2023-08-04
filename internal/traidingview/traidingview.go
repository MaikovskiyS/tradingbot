package traidingview

import (
	"fmt"
	"net/http"
	"time"
	"trrader/internal/traidingview/adx"
	"trrader/internal/traidingview/kaufman"
	"trrader/internal/traidingview/macd"
	"trrader/internal/traidingview/maintrend"
	"trrader/internal/traidingview/trend"
	"trrader/internal/traidingview/volatilityos"
	"trrader/internal/traidingview/volumeos"
)

const (
	Symbol = "OGNUSDT"
)

type Trader interface {
	PumpAndDump()
	ManagePosition(symbol string)
	StartTraiding(ticker, side string) error
}
type TraidingView struct {
	Trader     Trader
	Router     *http.ServeMux
	MainTrend  *maintrend.Indicator
	Trend      *trend.Indicator
	Volatility *volatilityos.Indicator
	Volume     *volumeos.Indicator
	Adx        *adx.Indicator
	Macd       *macd.Indicator
	Kaufman    *kaufman.Indicator
}

func New(svc Trader) *TraidingView {
	macd := macd.New()
	kaufman := kaufman.New()
	volatility := volatilityos.New()
	volume := volumeos.New()
	adx := adx.New()
	trend := trend.New()
	maintrend := maintrend.New()
	return &TraidingView{
		Trader:     svc,
		Router:     http.NewServeMux(),
		MainTrend:  maintrend,
		Trend:      trend,
		Volatility: volatility,
		Volume:     volume,
		Adx:        adx,
		Macd:       macd,
		Kaufman:    kaufman,
	}
}
func (t *TraidingView) RegisterRoutes() {

	t.Router.HandleFunc("/macd", t.Macd.GetAlert)
	t.Router.HandleFunc("/kaufman", t.Kaufman.GetAlert)
	t.Router.HandleFunc("/volatility", t.Volatility.GetAlert)
	t.Router.HandleFunc("/volume", t.Volume.GetAlert)
	t.Router.HandleFunc("/adx", t.Adx.GetAlert)
	t.Router.HandleFunc("/trend", t.Trend.GetAlert)
	t.Router.HandleFunc("/maintrend", t.MainTrend.GetAlert)

	//t.Router.HandleFunc("/nadaria", t.Nadaria.GetAlert)
}
func (t *TraidingView) Start() {
	//strategyAlert := StrategyAlert{}
	//ch := make(chan Alert)
	fmt.Println("Start listen alerts")
	go func() {
		fmt.Println("Start checking positions")
		t.Trader.ManagePosition(Symbol)
	}()
	// go func() {
	// 	t.Trader.PumpAndDump()
	// }()
	// for data := range t.Adx.Channel {
	// 	adx := SideAlert{Side: data.Side, Ticker: data.Ticker}
	// 	strategyAlert.Adx = adx
	// 	time.Sleep(1 * time.Second)

	// 	select {
	// 	case data := <-t.Volatility.Channel:
	// 		volatility := SideAlert{Side: data.Side, Ticker: data.Ticker}
	// 		fmt.Println("volatility chan case")
	// 		strategyAlert.Volatility = volatility
	// 		fmt.Println("set alert in volatility case")
	// 		time.Sleep(1 * time.Second)
	// 	default:
	// 		fmt.Println("no data from volatility")
	// 		continue
	// 	}
	// 	select {
	// 	case data := <-t.Volume.Channel:
	// 		volume := ActionAlert{Action: data.Action, Ticker: data.Ticker}
	// 		fmt.Println("volume chan case")
	// 		strategyAlert.Volume = volume
	// 		fmt.Println("set alert in volume case")
	// 		time.Sleep(1 * time.Second)
	// 	default:
	// 		fmt.Println("no data from volume")
	// 		continue
	// 	}
	// 	if strategyAlert.Adx.Ticker == "" && strategyAlert.Volatility.Ticker == "" && strategyAlert.Volume.Action == "" {
	// 		continue
	// 	} else {
	// 		if strategyAlert.Adx.Side == strategyAlert.Volatility.Side || strategyAlert.Adx.Ticker == strategyAlert.Volatility.Ticker || strategyAlert.Adx.Ticker == strategyAlert.Volume.Ticker {
	// 			err := t.Trader.StartTraiding(strategyAlert)
	// 			fmt.Println(err)
	// 			continue
	// 		}
	// 	}
	// }
	t.BreakOutStrategy()
}

// BreakOutStrategy
func (t *TraidingView) BreakOutStrategy() {
	alert := BreakOutAlert{}
	for {
		fmt.Println("waiting alerts")
		select {
		case data := <-t.MainTrend.Channel:
			fmt.Println("maintrend alert")
			maintrend := SideAlert{Side: data.Side, Ticker: data.Ticker, TimeStamp: data.TimeStamp}
			alert.Trend = maintrend
		case data := <-t.Trend.Channel:
			fmt.Println("trend alert")
			trend := SideAlert{Side: data.Side, Ticker: data.Ticker, TimeStamp: data.TimeStamp}
			alert.Trend = trend
		case data := <-t.Volatility.Channel:
			fmt.Println("volatility alert")
			volatility := ActionAlert{Action: data.Action, Ticker: data.Ticker, TimeStamp: data.TimeStamp}
			alert.Volatility = volatility
		case data := <-t.Volume.Channel:
			fmt.Println("volume alert")
			volume := ActionAlert{Action: data.Action, Ticker: data.Ticker, TimeStamp: data.TimeStamp}
			alert.Volume = volume
		}
		time.Sleep(1 * time.Second)
		fmt.Println(alert)
		err := alert.Validate()
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			fmt.Println("start traiding")
			t.Trader.StartTraiding(alert.Trend.Ticker, alert.Trend.Side)
		}

		time.Sleep(5 * time.Second)
	}

}

// ADX Strategy
func (t *TraidingView) AdxStrategy() {
	alert := AdxStrategyAlert{}
	for {
		fmt.Println("waiting alerts")
		select {
		case data := <-t.Adx.Channel:
			fmt.Println("adx alert")
			adx := SideAlert{Side: data.Side, Ticker: data.Ticker, TimeStamp: data.TimeStamp}
			alert.Adx = adx
		case data := <-t.Volatility.Channel:
			fmt.Println("volatility alert")
			volatility := SideAlert{Side: data.Side, Ticker: data.Ticker, TimeStamp: data.TimeStamp}
			alert.Volatility = volatility
		case data := <-t.Volume.Channel:
			fmt.Println("volume alert")
			volume := ActionAlert{Action: data.Action, Ticker: data.Ticker}
			alert.Volume = volume
		}
		time.Sleep(1 * time.Second)
		fmt.Println(alert)
		err := alert.Validate()
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			fmt.Println("start traiding")
			t.Trader.StartTraiding(alert.Adx.Ticker, alert.Adx.Side)
		}

		time.Sleep(5 * time.Second)
	}
}
