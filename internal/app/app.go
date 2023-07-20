package app

import (
	"trrader/internal/adapter"
	"trrader/internal/domain/service"
	"trrader/internal/server"
	"trrader/internal/traidingview"
)

func Run() error {
	tv := traidingview.New()
	//tv.GetData()
	rest := adapter.NewRestClient("url string", "apiKey string", "apiSecret string")
	tv.RegisterRoutes()
	tv.Start()
	svc := service.New(tv, rest)
	svc.StartTraiding()
	server := server.New(tv)
	server.Run()
	return nil
}
