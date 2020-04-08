package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/COVID-19-electronic-health-system/Coronalert/middleware"
	"github.com/COVID-19-electronic-health-system/Coronalert/router"
)

func main() {

	go middleware.StartPolling()

	r := router.Router()
	fmt.Println("server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
