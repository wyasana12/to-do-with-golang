package routes

import (
	"to-do-list-go/controllers"
	"to-do-list-go/middleware"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	router := r.PathPrefix("/user").Subrouter()

	router.Use(middleware.Auth)

	router.HandleFunc("/profile", controllers.Profile).Methods("GET")
}
