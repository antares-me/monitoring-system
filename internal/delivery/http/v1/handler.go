package v1

import (
	"antares-me/monitoring-system/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	monitoringService service.MonitoringService
}

func NewHandler(monService service.MonitoringService) *Handler {
	return &Handler{
		monitoringService: monService,
	}
}

func (h *Handler) Init(m *mux.Router) {
	h.MonitoringInit(m)
	h.StaticInit(m)
}
