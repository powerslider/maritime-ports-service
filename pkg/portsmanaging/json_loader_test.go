package portsmanaging_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/powerslider/maritime-ports-service/pkg/portsmanaging"
	"github.com/powerslider/maritime-ports-service/pkg/storage/memory"
)

var expectedPorts = map[string]*portsmanaging.MaritimePort{
	"AEAJM": {
		ID:          "AEAJM",
		Name:        "Ajman",
		City:        "Ajman",
		Country:     "United Arab Emirates",
		Alias:       []string{},
		Regions:     []string{},
		Coordinates: []float64{55.5136433, 25.4052165},
		Province:    "Ajman",
		Timezone:    "Asia/Dubai",
		Unlocs:      []string{"AEAJM"},
		Code:        "52000",
	},
	"AEAUH": {
		ID:          "AEAUH",
		Name:        "Abu Dhabi",
		Coordinates: []float64{54.37, 24.47},
		City:        "Abu Dhabi",
		Province:    "Abu Z¸aby [Abu Dhabi]",
		Country:     "United Arab Emirates",
		Alias:       []string{},
		Regions:     []string{},
		Timezone:    "Asia/Dubai",
		Unlocs:      []string{"AEAUH"},
		Code:        "52001",
	},
	"AEDXB": {
		ID:          "AEDXB",
		Name:        "Dubai",
		Coordinates: []float64{55.27, 25.25},
		City:        "Dubai",
		Province:    "Dubayy [Dubai]",
		Country:     "United Arab Emirates",
		Alias:       []string{},
		Regions:     []string{},
		Timezone:    "Asia/Dubai",
		Unlocs:      []string{"AEDXB"},
		Code:        "52005",
	},
}

func TestJSONLoader(t *testing.T) {
	portsStore := memory.NewPortsRepository()
	loader := portsmanaging.NewJSONLoader(portsStore)

	err := loader.LoadJSONFile("../../testdata/test_data_ports.json")
	require.NoError(t, err)

	storedPorts, err := portsStore.GetAllPorts()
	require.NoError(t, err)

	for _, port := range storedPorts {
		assert.Equal(t, port, expectedPorts[port.ID])
	}
}
