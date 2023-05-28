package handlers

import (
	"fmt"

	"github.com/gorilla/mux"
	_ "github.com/powerslider/maritime-ports-service/docs"
	"github.com/powerslider/maritime-ports-service/pkg/configs"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	// EndpointCreateOrUpdatePort is an HTTP endpoint for create or update port operation.
	EndpointCreateOrUpdatePort = "/api/v1/ports"
	// EndpointGetAllPorts is an HTTP endpoint for getting all ports operation.
	EndpointGetAllPorts = "/api/v1/ports"
	// EndpointGetPortByID is an HTTP endpoint for getting a port by ID operation.
	EndpointGetPortByID = "/api/v1/ports/{id}"
)

func registerHTTPRoutes(
	config *configs.Config, muxer *mux.Router, handler *PortsHandler) *mux.Router {
	muxer.HandleFunc(
		EndpointCreateOrUpdatePort,
		handler.CreateOrUpdatePort()).Methods("POST")
	muxer.HandleFunc(
		EndpointGetPortByID,
		handler.GetPort()).Methods("GET")
	muxer.HandleFunc(
		EndpointGetAllPorts,
		handler.GetAllPorts()).Methods("GET")

	swaggerJsonURL := fmt.Sprintf("http://%s:%d/swagger/doc.json", config.Host, config.Port)

	muxer.PathPrefix("/swagger/").
		Handler(httpSwagger.Handler(httpSwagger.URL(swaggerJsonURL))) // The url pointing to API definition

	return muxer
}
