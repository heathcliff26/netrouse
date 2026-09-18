package v1

//	@title			NetRouse API
//	@version		1.0
//	@description	Manage known hosts and send magic packets.

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1
//	@accept		json
//	@produce	json

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"log/slog"
	"net/http"

	"github.com/heathcliff26/netrouse/pkg/ping"
	"github.com/heathcliff26/netrouse/pkg/server/storage"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/heathcliff26/netrouse/pkg/utils"
	"github.com/heathcliff26/netrouse/pkg/wol"
)

type Response struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type apiHandler struct {
	storage *storage.Storage
}

func NewRouter(storage *storage.Storage) *http.ServeMux {
	handler := &apiHandler{
		storage: storage,
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /wake/{macAddr}", WakeHandler)
	router.HandleFunc("GET /hosts", handler.GetHostsHandler)
	router.HandleFunc("PUT /hosts", handler.AddHostHandler)
	router.HandleFunc("DELETE /hosts/{macAddr}", handler.RemoveHostHandler)
	router.HandleFunc("GET /hosts/status", handler.HostStatusHandler)
	return router
}

// @Summary		Wake up host
// @Description	Send a magic packet to the specified MAC address
//
// @Produce		json
// @Param			macAddr	path		string		true	"MAC address of the host"
// @Success		200		{object}	Response	"ok"
// @Failure		400		{object}	Response	"Invalid MAC address"
// @Failure		500		{object}	Response	"Failed to send magic packet"
// @Router			/wake/{macAddr} [get]
func WakeHandler(res http.ResponseWriter, req *http.Request) {
	macAddr := req.PathValue("macAddr")

	packet, err := wol.CreatePacket(macAddr)
	if err != nil {
		slog.Info("Client sent invalid MAC address", slog.String("mac", macAddr), slog.Any("error", err))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Invalid MAC address")
		return
	}

	err = packet.Send("")
	if err != nil {
		slog.Info("Failed to send magic packet", slog.String("mac", macAddr), slog.Any("error", err))
		res.WriteHeader(http.StatusInternalServerError)
		sendResponse(res, "Failed to send magic packet")
		return
	}

	slog.Info("Sent magic packet", slog.String("mac", macAddr))
	sendResponse(res, "")
}

// @Summary		Get hosts
// @Description	Fetch all known hosts
//
// @Produce		json
// @Success		200	{object}	[]types.Host	"List of all known hosts"
// @Failure		500	{object}	Response		"Failed to retrieve hosts from storage"
// @Router			/hosts [get]
func (h *apiHandler) GetHostsHandler(res http.ResponseWriter, req *http.Request) {
	hosts, err := h.storage.GetHosts()
	if err != nil {
		slog.Error("Failed to fetch hosts", "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		sendResponse(res, "Failed to fetch hosts")
		return
	}

	sendJSONResponse(res, hosts)
}

// @Summary		Add new host
// @Description	Add a new host to the known hosts
//
// @Accept			json
// @Produce		json
// @Param			payload	body		types.Host	true	"New host to add"
// @Success		200		{object}	Response	"ok"
// @Failure		400		{object}	Response	"Invalid MAC address or hostname"
// @Failure		403		{object}	Response	"Storage is readonly"
// @Failure		500		{object}	Response	"Failed to add host"
// @Router			/hosts [put]
func (h *apiHandler) AddHostHandler(res http.ResponseWriter, req *http.Request) {
	var host types.Host
	err := json.UnmarshalRead(req.Body, &host)
	if err != nil {
		slog.Debug("Client sent invalid host json", "error", err)
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Request body must be a valid host JSON object")
		return
	}

	if h.storage.Readonly() {
		slog.Debug("Client tried to add host while storage is readonly")
		res.WriteHeader(http.StatusForbidden)
		sendResponse(res, "Storage is readonly")
		return
	}

	if !utils.ValidateMACAddress(host.MAC) {
		slog.Debug("Client send invalid MAC address", slog.String("mac", host.MAC))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Invalid MAC address")
		return
	}

	if !utils.ValidateHostname(host.Name) {
		slog.Debug("Client send invalid hostname", slog.String("name", host.Name))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Invalid hostname")
		return
	}

	err = h.storage.AddHost(host)
	if err != nil {
		slog.Error("Failed to add host", "host", host, "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		sendResponse(res, "Failed to add host")
		return
	}

	slog.Info("Added host", "host", host)
	sendResponse(res, "")
}

// @Summary		Remove host
// @Description	Remove a host from the list of known hosts
//
// @Produce		json
// @Param			macAddr	path		string		true	"MAC address of the host"
// @Success		200		{object}	Response	"ok"
// @Failure		400		{object}	Response	"Invalid MAC address"
// @Failure		403		{object}	Response	"Storage is readonly"
// @Failure		500		{object}	Response	"Failed to remove host"
// @Router			/hosts/{macAddr} [delete]
func (h *apiHandler) RemoveHostHandler(res http.ResponseWriter, req *http.Request) {
	macAddr := req.PathValue("macAddr")

	if h.storage.Readonly() {
		slog.Debug("Client tried to remove host while storage is readonly")
		res.WriteHeader(http.StatusForbidden)
		sendResponse(res, "Storage is readonly")
		return
	}

	if !utils.ValidateMACAddress(macAddr) {
		slog.Debug("Client send invalid MAC address", slog.String("mac", macAddr))
		res.WriteHeader(http.StatusBadRequest)
		sendResponse(res, "Invalid MAC address")
		return
	}

	err := h.storage.RemoveHost(macAddr)
	if err != nil {
		slog.Error("Failed to remove host", "mac", macAddr, "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		sendResponse(res, "Failed to remove host")
		return
	}

	slog.Info("Removed host", slog.String("mac", macAddr))
	sendResponse(res, "")
}

// @Summary		Get host status
// @Description	Get the status of all hosts with an address to check if they are online
//
// @Produce		json
// @Success		200	{object}	[]types.HostStatus	"List of host statuses"
// @Failure		500	{object}	Response			"Failed to fetch or ping hosts"
// @Router			/hosts/status [get]
func (h *apiHandler) HostStatusHandler(res http.ResponseWriter, req *http.Request) {
	hosts, err := h.storage.GetHosts()
	if err != nil {
		slog.Error("Failed to fetch hosts", "error", err)
		res.WriteHeader(http.StatusInternalServerError)
		sendResponse(res, "Failed to fetch hosts")
		return
	}

	hostsToCheck := make([]types.Host, 0, len(hosts))
	for _, host := range hosts {
		if host.Address != "" {
			hostsToCheck = append(hostsToCheck, host)
		}
	}

	sendJSONResponse(res, ping.PingHosts(hostsToCheck))
}

func sendResponse(rw http.ResponseWriter, reason string) {
	response := Response{
		Status: "error",
		Reason: reason,
	}
	if reason == "" {
		response.Status = "ok"
	}

	sendJSONResponse(rw, response)
}

// Send an arbitrary JSON Object to the client
func sendJSONResponse(rw http.ResponseWriter, data any) {
	b, err := json.Marshal(data, jsontext.WithIndent("  "))
	if err != nil {
		slog.Error("Failed to create Response", "err", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	_, err = rw.Write(b)
	if err != nil {
		slog.Error("Failed to send response to client", "err", err)
	}
}
