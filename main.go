package main

import (
	"fmt"
	"log"
	"net/http"

	"./middleware"
	"./router"
)

func main() {

	go middleware.StartPolling()

	r := router.Router()
	fmt.Println("server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
