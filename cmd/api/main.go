package main

import (
	"log"
)

func main() {

	err := setupRouter().Run() // listens on 0.0.0.0:8080 by default
	if err != nil {
		log.Fatal(err)
	}
}
