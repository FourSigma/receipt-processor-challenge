package main

import "log"

func main() {
	api := NewAPI()
	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
