package portsmanaging

// PortsStore is a port interface representing operations on portsmanaging.MaritimePort entity.
type PortsStore interface {
	// UpsertPort inserts or modifies a new/existing portsmanaging.MaritimePort entity.
	UpsertPort(port *MaritimePort) (*MaritimePort, bool, error)

	// GetAllPorts returns all available ports from type portsmanaging.MaritimePort.
	GetAllPorts() ([]*MaritimePort, error)

	// GetPortByID returns n portsmanaging.MaritimePort identified by an available ID.
	GetPortByID(id string) (*MaritimePort, error)
}
