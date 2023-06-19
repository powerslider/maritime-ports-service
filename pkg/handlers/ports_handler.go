package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"


	"github.com/powerslider/maritime-ports-service/pkg/portsmanaging"

	"github.com/gorilla/mux"
	pkgErrors "github.com/pkg/errors"
)

// PortsService is a port interface for operations on portsmanaging.MaritimePort.
type PortsService interface {
	GetAllPorts() ([]*portsmanaging.MaritimePort, error)
	GetPortByID(ID string) (*portsmanaging.MaritimePort, error)
	CreateOrUpdatePort(p *portsmanaging.MaritimePort) (*portsmanaging.MaritimePort, bool, error)
}

// PortsHandler represents an HTTP handler for Ethereum block operations.
type PortsHandler struct {
	Service PortsService
}

// NewPortsHandler initializes a new instance of PortsHandler.
func NewPortsHandler(service PortsService) *PortsHandler {
	return &PortsHandler{
		Service: service,
	}
}

// GetAllPorts godoc
// @Summary Get all ports stored in the system.
// @Description Get all ports stored in the system.
// @Tags ports
// @Accept  json
// @Produce  json
// @Router /api/v1/ports [get]
func (h *PortsHandler) GetAllPorts() http.HandlerFunc {
	type response struct {
		Result []*portsmanaging.MaritimePort `json:"result"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		ports, err := h.Service.GetAllPorts()
		if err != nil {
			badRequestError(
				rw,
				pkgErrors.Wrapf(err, "could not get all ports"),
			)

			return
		}

		handleResponse(rw, response{
			Result: ports,
		})
	}
}

// GetPort godoc
// @Summary Get an existing port by ID.
// @Description Get an existing port by ID.
// @Tags ports
// @Accept  json
// @Produce  json
// @Param id path string true "MaritimePort ID"
// @Router /api/v1/ports/{id} [get]
func (h *PortsHandler) GetPort() http.HandlerFunc {
	type response struct {
		Result *portsmanaging.MaritimePort `json:"result"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, ok := vars["id"]
		if !ok {
			badRequestError(
				rw,
				errors.New("required path param 'id' is missing"),
			)

			return
		}

		p, err := h.Service.GetPortByID(id)
		if err != nil {
			badRequestError(
				rw,
				pkgErrors.Wrapf(err, "error getting port entry with ID '%s'", id),
			)

			return
		}

		if p == nil {
			notFoundError(
				rw,
				fmt.Errorf("port entry with ID '%s' not found", id),
			)

			return
		}

		handleResponse(rw, response{
			Result: p,
		})
	}
}

// CreateOrUpdatePort godoc
// @Summary Create a new port or update an existing one.
// @Description Create a new port or update an existing one.
// @Tags ports
// @Accept  json
// @Produce  json
// @Param request body portsmanaging.MaritimePort true "MaritimePort Entry"
// @Router /api/v1/ports [post]
func (h *PortsHandler) CreateOrUpdatePort() http.HandlerFunc {
	type response struct {
		Success bool   `json:"success"`
		Exists  bool   `json:"exists"`
		PortID  string `json:"port_id"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		var reqBody portsmanaging.MaritimePort

		reqBytes, errReqBytes := io.ReadAll(r.Body)
		errReqUnmarshal := json.Unmarshal(reqBytes, &reqBody)

		errReq := errors.Join(errReqBytes, errReqUnmarshal)
		if errReq != nil {
			badRequestError(
				rw,
				pkgErrors.Wrap(errReq, "could not unmarshal request params"),
			)

			return
		}

		if reqBody.ID == "" {
			badRequestError(
				rw,
				errors.New("required body param 'id' is missing"),
			)

			return
		}

		p, exists, err := h.Service.CreateOrUpdatePort(&reqBody)
		if err != nil {
			badRequestError(
				rw,
				pkgErrors.Wrap(errReq, "could not create/update port"),
			)

			return
		}

		handleResponse(rw, response{
			Success: true,
			Exists:  exists,
			PortID:  p.ID,
		})
	}
}

func handleResponse(rw http.ResponseWriter, resp any) {
	jsonResp, errRespMarshal := json.Marshal(resp)
	_, errRespWrite := rw.Write(jsonResp)

	errResp := errors.Join(errRespMarshal, errRespWrite)
	if errResp != nil {
		http.Error(rw, errResp.Error(), http.StatusInternalServerError)
	}
}

func badRequestError(rw http.ResponseWriter, err error) {
	errBytes, err := json.Marshal(struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{
		Status: http.StatusBadRequest,
		Error:  err.Error(),
	})

	if err == nil {
		http.Error(rw, string(errBytes), http.StatusBadRequest)
	} else {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func notFoundError(rw http.ResponseWriter, err error) {
	errBytes, err := json.Marshal(struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{
		Status: http.StatusNotFound,
		Error:  err.Error(),
	})

	if err == nil {
		http.Error(rw, string(errBytes), http.StatusNotFound)
	} else {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
