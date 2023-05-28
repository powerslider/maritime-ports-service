package handlers_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"

	"github.com/powerslider/maritime-ports-service/pkg/handlers"
	"github.com/powerslider/maritime-ports-service/pkg/portsmanaging"

	"github.com/powerslider/maritime-ports-service/pkg/storage/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tests = []struct {
	testCaseName             string
	httpMethod               string
	httpEndpoint             string
	httpPathParams           map[string]string
	httpRequestBody          string
	handlerFunc              func(portsHandler *handlers.PortsHandler) http.HandlerFunc
	expectedResponse         string
	expectedResponseFileName string
}{
	{
		testCaseName: "should return a correct response for getting all ports",
		httpMethod:   "GET",
		httpEndpoint: handlers.EndpointGetAllPorts,
		handlerFunc: func(portsHandler *handlers.PortsHandler) http.HandlerFunc {
			return portsHandler.GetAllPorts()
		},
		expectedResponseFileName: "get_all_ports_expected_response",
	},
	{
		testCaseName: "should return a correct response for creating a new port",
		httpMethod:   "POST",
		httpEndpoint: handlers.EndpointCreateOrUpdatePort,
		httpRequestBody: `
		{
			"id": "NEWPORT",
			"name": "Newest Port",
			"coordinates": [
			  123.321,
			  43.34
			],
			"city": "Some City",
			"country": "Some Country",
			"alias": [],
			"regions": [],
			"timezone": "My/Timezone",
			"unlocs": [
			  "NEWPORT"
			]
		}`,
		handlerFunc: func(portsHandler *handlers.PortsHandler) http.HandlerFunc {
			return portsHandler.CreateOrUpdatePort()
		},
		expectedResponse: `
		{
			"success": true,
			"exists": false,
			"port_id": "NEWPORT"
		}`,
	},
	{
		testCaseName: "should return a correct response for updating an existing port",
		httpMethod:   "POST",
		httpEndpoint: handlers.EndpointCreateOrUpdatePort,
		handlerFunc: func(portsHandler *handlers.PortsHandler) http.HandlerFunc {
			return portsHandler.CreateOrUpdatePort()
		},
		httpRequestBody: `
		{
			"id": "AEAJM",
			"city": "London",
			"country": "United Kingdom"
		}`,
		expectedResponse: `
		{
			"success": true,
			"exists": true,
			"port_id": "AEAJM"
		}`,
	},
	{
		testCaseName: "should return a correct response for querying an existing port",
		httpMethod:   "GET",
		httpEndpoint: handlers.EndpointGetPortByID,
		httpPathParams: map[string]string{
			"id": "AEDXB",
		},
		handlerFunc: func(portsHandler *handlers.PortsHandler) http.HandlerFunc {
			return portsHandler.GetPort()
		},
		expectedResponse: `
		{
		   "result":{
			  "id":"AEDXB",
			  "name":"Dubai",
			  "city":"Dubai",
			  "country":"United Arab Emirates",
			  "alias":[],
			  "regions":[],
			  "coordinates":[
				 55.27,
				 25.25
			  ],
			  "province":"Dubayy [Dubai]",
			  "timezone":"Asia/Dubai",
			  "unlocs":[
				 "AEDXB"
			  ],
			  "code":"52005"
		   }
		}`,
	},
}

func TestPortsHandlerCorrectResponses(t *testing.T) {
	for _, test := range tests {
		portsHandler := setupHandler(t)

		var reqBody io.Reader

		if len(test.httpRequestBody) > 0 {
			reqBody = bytes.NewBuffer([]byte(test.httpRequestBody))
		}

		req, errReq := http.NewRequest(test.httpMethod, test.httpEndpoint, reqBody)
		require.NoError(t, errReq)

		rr := httptest.NewRecorder()

		if len(test.httpPathParams) > 0 {
			req = mux.SetURLVars(req, test.httpPathParams)
		}

		handler := test.handlerFunc(portsHandler)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		verifyExpectedResponse(t, test.expectedResponse, test.expectedResponseFileName, rr)
	}
}

func setupHandler(t *testing.T) *handlers.PortsHandler {
	portsStore := memory.NewPortsRepository()
	portsService := portsmanaging.NewService(portsStore)
	portsHandler := handlers.NewPortsHandler(portsService)
	loader := portsmanaging.NewJSONLoader(portsStore)

	err := loader.LoadJSONFile("../../testdata/test_data_ports.json")
	require.NoError(t, err)

	return portsHandler
}

func verifyExpectedResponse(
	t *testing.T,
	expectedResponse string,
	expectedResponseFileName string,
	respRec *httptest.ResponseRecorder,
) {
	var expected []byte

	if len(expectedResponseFileName) > 0 {
		filePath, errFilePath := filepath.Abs(fmt.Sprintf("../../testdata/%s.json", expectedResponseFileName))

		expResp, errFile := os.ReadFile(filePath)
		if err := errors.Join(errFilePath, errFile); err != nil {
			t.Errorf("could not read expected response JSON file: %v", err)
		}

		expected = expResp
	} else {
		expected = []byte(expectedResponse)
	}

	assert.JSONEq(t, string(expected), respRec.Body.String())
}
