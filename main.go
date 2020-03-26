package main

import (
	"fmt"
	"log"
	"net/http"

	"./router"
)

func main() {

	r := router.Router()
	fmt.Println("server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
