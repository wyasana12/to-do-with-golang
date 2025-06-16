package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RouteIndex(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	AuthRoutes(api)
	UserRoutes(api)
	TodoRoutes(api)

	api.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api.yaml")
	}).Methods("GET")
}
