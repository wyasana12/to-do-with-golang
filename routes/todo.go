package routes

import (
	todocontroller "to-do-list-go/controllers/todoController"
	"to-do-list-go/middleware"

	"github.com/gorilla/mux"
)

func TodoRoutes(r *mux.Router) {
	router := r.PathPrefix("/todos").Subrouter()
	router.Use(middleware.Auth)
	router.HandleFunc("", todocontroller.Index).Methods("GET")
	router.HandleFunc("", todocontroller.Create).Methods("POST")
	router.HandleFunc("/bulk-delete", todocontroller.BulkDestroy).Methods("DELETE")

	AttachmentRoutes(router)
	TrashRoutes(router)

	router.HandleFunc("/{id}", todocontroller.Detail).Methods("GET")
	router.HandleFunc("/{id}", todocontroller.Update).Methods("PUT")
	router.HandleFunc("/{id}", todocontroller.Destroy).Methods("DELETE")
}
