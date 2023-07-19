package app

import (
	"fmt"
	"trrader/internal/adapter"
	"trrader/internal/domain/service"
	"trrader/internal/traidingview"
)

func Run() error {
	adapter.HeyfromAdapter()
	service.Hellosvc()
	traidingview.HH()
	fmt.Println("servicestarted")
	return nil
}
