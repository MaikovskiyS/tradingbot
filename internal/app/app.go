package app

import (
	"os"
	"trrader/internal/adapter"
	"trrader/internal/adapter/bybit"
	"trrader/internal/domain/service"
	"trrader/internal/server"
	"trrader/internal/traidingview"
)

func Run() error {
	secret := os.Getenv("secretkey")
	api := os.Getenv("apikey")
	tv := traidingview.New()
	//tv.GetData()
	rest := adapter.NewRestClient(bybit.BytickMainnetBaseURL, api, secret)
	tv.RegisterRoutes()
	tv.Start()
	svc := service.New(tv, rest)
	svc.StartTraiding()
	server := server.New(tv)
	server.Run()
	return nil
}
