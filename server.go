package goapi

import (
	"fmt"
	"log"
	"net/http"
)

type ServiceInterface interface {
	Controller(config *Config, w http.ResponseWriter, req *http.Request)
}

var Service ServiceInterface

// StartService ...
func StartService(srv ServiceInterface) {
	Service = srv

	// Read configuration from the file path
	config := ReadConfig()

	// Create a new Spanner client
	config.SpannerClient()
	defer SpannerClient.Close()

	http.HandleFunc(config.Service.Path, config.Router) // Load all the routes

	fmt.Println("Starting service:", config)

	addr := config.Service.Host + ":" + config.Service.Port
	err := http.ListenAndServe(addr, nil) // Start the API server
	if err != nil {
		log.Fatal("Error! server failed to start.", err)
	}
}
