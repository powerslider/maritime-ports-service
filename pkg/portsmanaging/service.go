package portsmanaging

// Service represents execution of business logic upon portsmanaging.MaritimePort.
type Service struct {
	Repository PortsStore
}

// NewService is a constructor function for Service.
func NewService(repository PortsStore) *Service {
	return &Service{
		Repository: repository,
	}
}

// GetAllPorts returns all ports of type portsmanaging.MaritimePort stored in the system.
func (h *Service) GetAllPorts() ([]*MaritimePort, error) {
	return h.Repository.GetAllPorts()
}

// GetPortByID returns a porn given a port ID.
func (h *Service) GetPortByID(ID string) (*MaritimePort, error) {
	return h.Repository.GetPortByID(ID)
}

// CreateOrUpdatePort add a new port entry of type portsmanaging.MaritimePort or updates an existing one.
func (h *Service) CreateOrUpdatePort(p *MaritimePort) (*MaritimePort, bool, error) {
	return h.Repository.UpsertPort(p)
}
