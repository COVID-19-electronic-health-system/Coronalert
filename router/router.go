package router

import (
	"../middleware"
	"github.com/gorilla/mux"
)

// Router exposes routes to main
func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/subscribe", middleware.Subscribe).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/unsubscribe", middleware.Unsubscribe).Methods("POST", "OPTIONS")

	return router
}
