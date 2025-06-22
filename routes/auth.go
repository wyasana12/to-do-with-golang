package routes

import (
	"to-do-list-go/controllers"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	router := r.PathPrefix("/auth").Subrouter()

	router.HandleFunc("/register", controllers.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/verify-email", controllers.VerifyEmail).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", controllers.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/forgot-password", controllers.RequestResetPassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/reset-password", controllers.ResetPassword).Methods("POST", "OPTIONS")
}
