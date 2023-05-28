package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	pkgErrors "github.com/pkg/errors"

	"github.com/powerslider/maritime-ports-service/pkg/entity"
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

// UpsertPort inserts a new entity.Port entity.
func (r *PortsRepository) UpsertPort(port *entity.Port) (*entity.Port, bool, error) {
	p, loaded := r.store.LoadOrStore(port.ID, port)

	if loaded {
		updatePortBytes, errMarshal := json.Marshal(port)
		errUnmarshal := json.Unmarshal(updatePortBytes, p)

		if err := errors.Join(errMarshal, errUnmarshal); err != nil {
			return nil, loaded, pkgErrors.Wrapf(
				err, "error: failed update of existing port with ID '%s'", port.ID)
		}

		updatedPort, ok := p.(*entity.Port)
		if !ok {
			return nil, loaded, fmt.Errorf("error: updated port entry is corrupt: %s", fmt.Sprint(p))
		}

		r.store.Store(updatedPort.ID, updatedPort)

		return updatedPort, loaded, nil
	}

	return p.(*entity.Port), loaded, nil
}

// GetAllPorts returns all available ports from type entity.Port.
func (r *PortsRepository) GetAllPorts() ([]*entity.Port, error) {
	var err error

	pp := make([]*entity.Port, 0)

	r.store.Range(func(key, value any) bool {
		p, ok := value.(*entity.Port)
		if ok {
			pp = append(pp, p)
		} else {
			err = fmt.Errorf("error: queried port data is corrupt: %s", fmt.Sprint(value))
		}

		return ok
	})

	return pp, err
}

// GetPortByID returns an entity.Port identified by an available ID.
func (r *PortsRepository) GetPortByID(id string) (*entity.Port, error) {
	v, loaded := r.store.Load(id)
	if loaded {
		p, ok := v.(*entity.Port)
		if ok {
			return p, nil
		}

		return nil, fmt.Errorf(
			"error: queried port data for entry with ID '%s' is corrupt: %s", id, fmt.Sprint(v))
	}

	return nil, nil
}
