package main

import (
	"log"
	"microservice/bootstrap"
)

// Hello
func main() {
	err := bootstrap.Run()
	if err != nil {
		log.Fatal(err)
	}
}
