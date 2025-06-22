package routes

import (
	trashcontroller "to-do-list-go/controllers/trashController"

	"github.com/gorilla/mux"
)

func TrashRoutes(r *mux.Router) {
	router := r.PathPrefix("/trash").Subrouter()

	router.HandleFunc("", trashcontroller.Trash).Methods("GET", "OPTIONS")
	router.HandleFunc("/bulk-restore", trashcontroller.BulkRestore).Methods("PUT", "OPTIONS")
	router.HandleFunc("/bulk-permanent-delete", trashcontroller.BulkDelete).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/restore/{id}", trashcontroller.Restore).Methods("PUT", "OPTIONS")
	router.HandleFunc("/permanent-delete/{id}", trashcontroller.PermanentDelete).Methods("DELETE", "OPTIONS")
}
