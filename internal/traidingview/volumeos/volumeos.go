package volumeos

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	fmt.Println("in volume get alert")
	err := json.NewDecoder(r.Body).Decode(&i.Data)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("alert:", i.Data)
	i.Channel <- i.Data
}
