package handlers

import (
	"github.com/gorilla/mux"
	"github.com/powerslider/maritime-ports-service/pkg/configs"
)

// InitializeHandlers registers HTTP routes and wires dependencies for HTTP handlers.
func InitializeHandlers(
	config *configs.Config,
	router *mux.Router,
	service PortsService,
) *mux.Router {
	handler := NewPortsHandler(service)

	registerHTTPRoutes(config, router, handler)

	return router
}
