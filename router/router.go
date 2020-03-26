package router

import (
	"../middleware"
	"github.com/gorilla/mux"
)

// Router exposes routes to main
func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/sendSMS", middleware.SendSMS).Methods("POST", "OPTIONS")

	return router
}
