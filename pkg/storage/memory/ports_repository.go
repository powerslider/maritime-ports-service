package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/powerslider/maritime-ports-service/pkg/portsmanaging"

	pkgErrors "github.com/pkg/errors"
)

// PortsRepository holds the CRUD db operations for CasinoRoundBet.
type PortsRepository struct {
	store sync.Map
}

// NewPortsRepository is a constructor function for PortsRepository.
func NewPortsRepository() *PortsRepository {
	return &PortsRepository{
		store: sync.Map{},
	}
}

// UpsertPort inserts or modifies a new/existing portsmanaging.MaritimePort entity.
func (r *PortsRepository) UpsertPort(port *portsmanaging.MaritimePort) (*portsmanaging.MaritimePort, bool, error) {
	p, loaded := r.store.LoadOrStore(port.ID, port)

	if loaded {
		updatePortBytes, errMarshal := json.Marshal(port)
		errUnmarshal := json.Unmarshal(updatePortBytes, p)

		if err := errors.Join(errMarshal, errUnmarshal); err != nil {
			return nil, loaded, pkgErrors.Wrapf(
				err, "error: failed update of existing port with ID '%s'", port.ID)
		}

		updatedPort, ok := p.(*portsmanaging.MaritimePort)
		if !ok {
			return nil, loaded, fmt.Errorf("error: updated port entry is corrupt: %s", fmt.Sprint(p))
		}

		r.store.Store(updatedPort.ID, updatedPort)

		return updatedPort, loaded, nil
	}

	return p.(*portsmanaging.MaritimePort), loaded, nil
}

// GetAllPorts returns all available ports from type portsmanaging.MaritimePort.
func (r *PortsRepository) GetAllPorts() ([]*portsmanaging.MaritimePort, error) {
	var err error

	pp := make([]*portsmanaging.MaritimePort, 0)

	r.store.Range(func(key, value any) bool {
		p, ok := value.(*portsmanaging.MaritimePort)
		if ok {
			pp = append(pp, p)
		} else {
			err = fmt.Errorf("error: queried port data is corrupt: %s", fmt.Sprint(value))
		}

		return ok
	})

	return pp, err
}

// GetPortByID returns n portsmanaging.MaritimePort identified by an available ID.
func (r *PortsRepository) GetPortByID(id string) (*portsmanaging.MaritimePort, error) {
	v, loaded := r.store.Load(id)
	if loaded {
		p, ok := v.(*portsmanaging.MaritimePort)
		if ok {
			return p, nil
		}

		return nil, fmt.Errorf(
			"error: queried port data for entry with ID '%s' is corrupt: %s", id, fmt.Sprint(v))
	}

	return nil, nil
}
