package routes

import "github.com/gorilla/mux"

func RouteIndex(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	AuthRoutes(api)
	UserRoutes(api)
	TodoRoutes(api)
}
