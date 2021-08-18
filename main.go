package main

import (
	"log"
	"wallet-engine/server"
)

func main() {
	if err := server.Init(); err != nil {
		log.Fatalln(err)
	}
}
