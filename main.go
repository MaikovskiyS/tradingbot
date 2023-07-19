package main

import (
	"fmt"
	"log"
	"trrader/internal/app"
)

func main() {
	fmt.Println("service starting")
	err := app.Run()
	if err != nil {
		log.Fatal()
	}

}
