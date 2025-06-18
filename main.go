package main

import (
	"fmt"
	"net/http"
	"to-do-list-go/config"
	"to-do-list-go/middleware"
	"to-do-list-go/routes"
	"to-do-list-go/schedular"
	"to-do-list-go/static"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()

	r := mux.NewRouter()

	r.Use(middleware.CORSmiddleware)

	routes.RouteIndex(r)

	schedular.StartNotificationSchedular()

	static.StartUploadFileServer(r, "./uploads")

	static.StartSwaggerUI(r, "./web/swagger-ui", "./web/api.yaml")
	// r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	path, _ := route.GetPathTemplate()
	// 	methods, _ := route.GetMethods()
	// 	log.Infof("Route registered: %s %v", path, methods)
	// 	return nil
	// })

	log.Println("Server Running On Port", config.ENV.PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", config.ENV.PORT), r)
}
