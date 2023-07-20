package kaufman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Indicator struct {
	Channel chan Data
	Data    Data
}

func New() *Indicator {
	return &Indicator{
		Channel: make(chan Data),
		Data:    Data{},
	}
}
func (i *Indicator) GetAlert(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&i.Data)
	if err != nil {
		fmt.Println(err)
	}
	for {
		i.Data.Ticker = "BTC"
		i.Data.Side = "Sell"
		i.Channel <- i.Data
		time.Sleep(time.Second * 30)
	}

}
func (i *Indicator) Testt() {
	for {

		i.Data.Ticker = "OPUSDT"
		i.Data.Side = "Sell"
		fmt.Println("data to chan in kaufman testt")
		i.Channel <- i.Data

		time.Sleep(time.Second * 10)
	}

}
