package portsmanaging

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	pkgErrors "github.com/pkg/errors"

	"github.com/powerslider/maritime-ports-service/pkg/entity"
)

// JSONLoader is a service responsible for loading json data.
type JSONLoader struct {
	Repository PortsStore
}

// NewJSONLoader is a constructor function for JSONLoader.
func NewJSONLoader(repository PortsStore) *JSONLoader {
	return &JSONLoader{
		Repository: repository,
	}
}

// LoadJSONFile reads a JSON file and delegates loading to Load method.
func (l *JSONLoader) LoadJSONFile(jsonFilePath string) error {
	dataFilePath, errPath := filepath.Abs(jsonFilePath)
	portsFixtures, errFile := os.Open(dataFilePath)

	if err := errors.Join(errPath, errFile); err != nil {
		return pkgErrors.Wrapf(err, "cannot access ports data from file %s", dataFilePath)
	}

	if err := l.Load(portsFixtures); err != nil {
		return pkgErrors.Wrapf(err, "cannot load ports from file: %s", dataFilePath)
	}

	return nil
}

// Load stores JSON data in chunks via PortsStore.
func (l *JSONLoader) Load(r io.Reader) error {
	dec := json.NewDecoder(r)

	var (
		token json.Token
		err   error
	)

	_, err = dec.Token()
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	for dec.More() {
		token, err = dec.Token()
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		var p entity.Port

		err = dec.Decode(&p)
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		p.ID = fmt.Sprint(token)

		_, _, err = l.Repository.UpsertPort(&p)
		if err != nil {
			return pkgErrors.WithStack(err)
		}
	}

	return nil
}
