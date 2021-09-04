package main

import (
	"log"

	"github.com/luisfernandogaido/cnc/ncw"
)

func main() {
	_, err := ncw.Simulacoes()
	if err != nil {
		log.Fatal(err)
	}
}
