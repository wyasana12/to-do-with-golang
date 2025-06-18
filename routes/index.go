package routes

import (
	"net/http"
	"to-do-list-go/middleware"

	"github.com/gorilla/mux"
)

func RouteIndex(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	r.Use(middleware.CORSmiddleware)

	AuthRoutes(api)
	UserRoutes(api)
	TodoRoutes(api)

	api.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api.yaml")
	}).Methods("GET")
}
