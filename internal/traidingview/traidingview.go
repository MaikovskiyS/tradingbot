package traidingview

import (
	"net/http"
	"sync"
	"trrader/internal/traidingview/kaufman"
	"trrader/internal/traidingview/macd"
)

// type TvGetter interface {
// 	GetData() Alert
// }
type TraidingView struct {
	wg      *sync.WaitGroup
	Router  *http.ServeMux
	Macd    *macd.Indicator
	Kaufman *kaufman.Indicator
	AlertCh chan Alert
}
type Alert struct {
	Ticker string
	Side   string
}

func New() *TraidingView {
	macd := macd.New()
	kaufman := kaufman.New()
	alertch := make(chan Alert)
	return &TraidingView{
		wg:      &sync.WaitGroup{},
		Router:  http.NewServeMux(),
		Macd:    macd,
		Kaufman: kaufman,
		AlertCh: alertch,
	}
}
func (t *TraidingView) RegisterRoutes() {

	t.Router.HandleFunc("/macd", t.Macd.GetAlert)
	t.Router.HandleFunc("kauman", t.Kaufman.GetAlert)
}
func (t *TraidingView) Start() {
	//t.wg.Add(1)
	go func() {
		//	defer t.wg.Done()
		t.Macd.Testt()
	}()
	//	t.wg.Add(1)
	go func() {
		//	defer t.wg.Done()
		t.Kaufman.Testt()
	}()
	//t.wg.Wait()
}

func (t *TraidingView) GetData() Alert {
	var alert Alert
	macd := <-t.Macd.Channel

	kaufman := <-t.Kaufman.Channel

	if macd.Side == kaufman.Side {

		alert.Ticker = macd.Ticker
		alert.Side = macd.Side
		//t.AlertCh <- alert
	}
	return alert
}
