package main

import (
	"log"

	"github.com/FourSigma/receipt-processor-challenge/pkg/api"
)

func main() {
	a := api.New()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
