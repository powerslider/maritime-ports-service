package portsmanaging

import (
	"github.com/powerslider/maritime-ports-service/pkg/entity"
)

// Service represents execution of business logic upon entity.Port.
type Service struct {
	Repository PortsStore
}

// NewService is a constructor function for Service.
func NewService(repository PortsStore) *Service {
	return &Service{
		Repository: repository,
	}
}

// GetAllPorts returns all ports of type entity.Port stored in the system.
func (h *Service) GetAllPorts() ([]*entity.Port, error) {
	return h.Repository.GetAllPorts()
}

// GetPortByID returns a porn given a port ID.
func (h *Service) GetPortByID(ID string) (*entity.Port, error) {
	return h.Repository.GetPortByID(ID)
}

// CreateOrUpdatePort add a new port entry of type entity.Port or updates an existing one.
func (h *Service) CreateOrUpdatePort(p *entity.Port) (*entity.Port, bool, error) {
	return h.Repository.UpsertPort(p)
}
