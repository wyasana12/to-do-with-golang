package routes

import (
	attachmentcontroller "to-do-list-go/controllers/attachmentController"

	"github.com/gorilla/mux"
)

func AttachmentRoutes(r *mux.Router) {
	router := r.PathPrefix("/{id}").Subrouter()

	router.HandleFunc("/attachment", attachmentcontroller.Upload).Methods("POST", "OPTIONS")
	router.HandleFunc("/attachment/{attachment_id}/download", attachmentcontroller.Download).Methods("GET", "OPTIONS")
	router.HandleFunc("/attachment/{attachment_id}", attachmentcontroller.Delete).Methods("DELETE", "OPTIONS")
}
