package traidingview

import (
	"fmt"
	"net/http"
	"time"
	"trrader/internal/traidingview/adx"
	"trrader/internal/traidingview/kaufman"
	"trrader/internal/traidingview/macd"
	"trrader/internal/traidingview/volatilityos"
	"trrader/internal/traidingview/volumeos"
)

const (
	Symbol = "OPUSDT"
)

type Trader interface {
	ManagePosition(symbol string)
	StartTraiding(a StrategyAlert) error
}
type TraidingView struct {
	Trader     Trader
	Router     *http.ServeMux
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
	return &TraidingView{
		Trader:     svc,
		Router:     http.NewServeMux(),
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
	t.AdxStrategy()
}
func (t *TraidingView) AdxStrategy() {
	alert := StrategyAlert{}
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
			t.Trader.StartTraiding(alert)
		}

		time.Sleep(5 * time.Second)
	}

}
