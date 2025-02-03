package main

import (
	"github.com/FourSigma/receipt-processor-challenge/pkg/api"
)

func main() {
	a := api.New()
	a.Run()
}
