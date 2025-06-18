package static

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func StartUploadFileServer(r *mux.Router, uploadDir string) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			log.Fatalf("failed to create upload directory %s: %v", uploadDir, err)
		}
		log.Infof("created upload directory: %s", uploadDir)
	}

	uploadHandler := http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir)))
	r.PathPrefix("/uploads/").Handler(uploadHandler)

	log.Infof("serving uploaded files from %s at /uploads/", uploadDir)
}

func StartSwaggerUI(r *mux.Router, swaggerUIDir string, apiYAMLPath string) {
	swaggerHandler := http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir(swaggerUIDir)))
	r.PathPrefix("/swagger-ui/").Handler(swaggerHandler)

	r.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, apiYAMLPath)
	}).Methods("GET")

	log.Infof("serving swagger UI from %s at /swagger-ui/ and API spec from %s at /swagger.yaml", swaggerUIDir, apiYAMLPath)
}
