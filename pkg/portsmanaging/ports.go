package portsmanaging

import "github.com/powerslider/maritime-ports-service/pkg/entity"

// PortsStore is a port interface representing operations on entity.Port entity.
type PortsStore interface {
	UpsertPort(port *entity.Port) (*entity.Port, bool, error)

	// GetAllPorts returns all available ports from type entity.Port.
	GetAllPorts() ([]*entity.Port, error)

	// GetPortByID returns a entity.Port identified by an available ID.
	GetPortByID(id string) (*entity.Port, error)
}
