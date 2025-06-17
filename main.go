package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"to-do-list-go/config"
	"to-do-list-go/routes"
	"to-do-list-go/schedular"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()

	r := mux.NewRouter()
	routes.RouteIndex(r)

	schedular.StartNotificationSchedular()
	// r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	path, _ := route.GetPathTemplate()
	// 	methods, _ := route.GetMethods()
	// 	log.Infof("Route registered: %s %v", path, methods)
	// 	return nil
	// })

	staticDir := "./web"

	swaggerUIHandler := http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir(filepath.Join(staticDir, "swagger-ui"))))
	r.PathPrefix("/swagger-ui/").Handler(swaggerUIHandler)

	r.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath.Join(staticDir, "api.yaml"))
	}).Methods("GET")

	log.Println("Server Running On Port", config.ENV.PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", config.ENV.PORT), r)
}
