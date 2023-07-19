package main

import (
	"fmt"
	"log"
	"trrader/internal/adapter"
	"trrader/internal/app"
	"trrader/internal/traidingview"
)

func main() {
	fmt.Println("hello second")
	err := app.Run()
	if err != nil {
		log.Fatal()
	}
	traidingview.HH()
	adapter.HeyfromAdapter()
}
