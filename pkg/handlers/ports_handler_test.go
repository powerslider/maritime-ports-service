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

	"github.com/kinbiko/jsonassert"

	"github.com/gorilla/mux"

	"github.com/powerslider/maritime-ports-service/pkg/handlers"
	"github.com/powerslider/maritime-ports-service/pkg/portsmanaging"

	"github.com/powerslider/maritime-ports-service/pkg/storage/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tests struct {
	testCaseName             string
	httpMethod               string
	httpEndpoint             string
	httpPathParams           map[string]string
	httpRequestBody          string
	handlerFunc              func(portsHandler *handlers.PortsHandler) http.HandlerFunc
	expectedResponse         string
	expectedResponseFileName string
	expectedResponseCode     int
}

func TestPortsHandlerCorrectResponses(t *testing.T) {
	var testData = []tests{
		{
			testCaseName: "should return a correct response for getting all ports",
			httpMethod:   "GET",
			httpEndpoint: handlers.EndpointGetAllPorts,
			handlerFunc: func(portsHandler *handlers.PortsHandler) http.HandlerFunc {
				return portsHandler.GetAllPorts()
			},
			expectedResponseCode:     http.StatusOK,
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
			expectedResponseCode: http.StatusOK,
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
			expectedResponseCode: http.StatusOK,
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
			expectedResponseCode: http.StatusOK,
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
		{
			testCaseName: "should return a correct response for non existent port",
			httpMethod:   "GET",
			httpEndpoint: handlers.EndpointGetPortByID,
			httpPathParams: map[string]string{
				"id": "NONEXISTENT",
			},
			handlerFunc: func(portsHandler *handlers.PortsHandler) http.HandlerFunc {
				return portsHandler.GetPort()
			},
			expectedResponseCode: http.StatusNotFound,
			expectedResponse: `
			{
			   "status": 404,
			   "error": "port entry with ID 'NONEXISTENT' not found"
			}`,
		},
	}

	ja := jsonassert.New(t)

	for _, test := range testData {
		t.Log(test.testCaseName)

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

		assert.Equal(t, test.expectedResponseCode, rr.Code)

		verifyExpectedResponse(t, ja, test.expectedResponse, test.expectedResponseFileName, rr)
	}
}

func TestPortsHandlerErroneousResponses(t *testing.T) {
	var testData = []tests{
		{
			testCaseName: "should return a validation error for a missing 'id' param " +
				"in the request body when creating a new port",
			httpMethod:   "POST",
			httpEndpoint: handlers.EndpointCreateOrUpdatePort,
			httpRequestBody: `
			{
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
			expectedResponseCode: http.StatusBadRequest,
			expectedResponse: `
			{
				"status": 400,
				"error": "required body param 'id' is missing"
			}`,
		},
		{
			testCaseName: "should return a request body unmarshal error when creating a new port",
			httpMethod:   "POST",
			httpEndpoint: handlers.EndpointCreateOrUpdatePort,
			httpRequestBody: `
			{
				"name": "Newest Port"
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
			expectedResponseCode: http.StatusBadRequest,
			expectedResponse: `
			{
				"status": 400,
				"error": "could not unmarshal request params: invalid character '\"' after object key:value pair"
			}`,
		},
		{
			testCaseName: "should return a validation error for a missing 'id' path param when querying a port by ID",
			httpMethod:   "GET",
			httpEndpoint: handlers.EndpointGetPortByID,
			handlerFunc: func(portsHandler *handlers.PortsHandler) http.HandlerFunc {
				return portsHandler.GetPort()
			},
			expectedResponseCode: http.StatusBadRequest,
			expectedResponse: `
			{
			   "status": 400,
			   "error": "required path param 'id' is missing"
			}`,
		},
	}

	ja := jsonassert.New(t)

	for _, test := range testData {
		t.Log(test.testCaseName)

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

		assert.Equal(t, test.expectedResponseCode, rr.Code)

		verifyExpectedResponse(t, ja, test.expectedResponse, test.expectedResponseFileName, rr)
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
	ja *jsonassert.Asserter,
	expectedResponse string,
	expectedResponseFileName string,
	respRec *httptest.ResponseRecorder,
) {
	var expected []byte

	if len(expectedResponseFileName) > 0 {
		filePath, errFilePath := filepath.Abs(fmt.Sprintf("../../testdata/%s.json", expectedResponseFileName))

		expResp, errFile := os.ReadFile(filePath)
		require.NoError(t, errors.Join(errFilePath, errFile))

		expected = expResp
	} else {
		expected = []byte(expectedResponse)
	}

	ja.Assertf(respRec.Body.String(), string(expected))
}
