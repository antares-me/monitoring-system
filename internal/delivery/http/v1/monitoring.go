package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) MonitoringInit(m *mux.Router) {
	m.HandleFunc("/api", h.handleConnection).Methods("GET")
}

func (h *Handler) handleConnection(w http.ResponseWriter, r *http.Request) {
	answer, err := json.Marshal(h.monitoringService.GetStatus(r.Context()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}
