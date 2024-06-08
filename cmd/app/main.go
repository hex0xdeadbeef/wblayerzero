package main

import (
	"log"
	"wblayerzero/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
