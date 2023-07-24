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

	bybit := adapter.NewRestClient(bybit.BytickMainnetBaseURL, api, secret)
	svc := service.New(bybit)
	tv := traidingview.New(svc)
	tv.RegisterRoutes()
	server := server.New(tv)
	server.Run()
	tv.Start()

	return nil
}
