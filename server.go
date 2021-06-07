package goapi

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
)

type ServiceInterface interface {
	Controller(config *Config, w http.ResponseWriter, req *http.Request)
}

var service ServiceInterface
var DB *spanner.Client

// StartService ...
func StartService(srv ServiceInterface) {
	service = srv

	// Read configuration from the file path
	config := ReadConfig()

	// Create a new Spanner client
	DB = config.newSpannerClient()
	defer DB.Close()

	http.HandleFunc(config.Service.Path, config.Router) // Load all the routes

	fmt.Println("Starting service:", config)

	addr := config.Service.Host + ":" + config.Service.Port
	err := http.ListenAndServe(addr, nil) // Start the API server
	if err != nil {
		log.Fatal("Error! server failed to start.", err)
	}
}
