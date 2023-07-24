package volatilityos

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
	fmt.Println("in volatilityos get alert")
	err := json.NewDecoder(r.Body).Decode(&i.Data)
	if err != nil {
		fmt.Println(err)
	}
	mark := time.Now().Unix()
	i.Data.TimeStamp = mark
	fmt.Println("volatility alert time:", i.Data.TimeStamp)
	i.Channel <- i.Data
}
