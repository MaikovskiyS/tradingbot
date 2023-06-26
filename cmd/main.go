package main

import (
	"fmt"
	"log"
	"traider/internal/app"
)

func main() {
	fmt.Println("hello second")
	err := app.Run()
	if err != nil {
		log.Fatal()
	}
}
