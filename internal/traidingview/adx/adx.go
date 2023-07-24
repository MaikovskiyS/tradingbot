package adx

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
	fmt.Println("in adx get alert")
	err := json.NewDecoder(r.Body).Decode(&i.Data)
	if err != nil {
		fmt.Println(err)
	}
	mark := time.Now().Unix()
	i.Data.TimeStamp = mark
	fmt.Println("adx alert time:", i.Data.TimeStamp)
	//fmt.Println("alert:", i.Data)
	i.Channel <- i.Data

}

// for local tests
func (i *Indicator) Testt() {
	for {
		i.Data.Ticker = "OPUSDT"
		i.Data.Side = "Sell"
		//	time.Sleep(time.Second * 10)
		fmt.Println("data to chan in kaufman testt")

		i.Channel <- i.Data

	}

}
