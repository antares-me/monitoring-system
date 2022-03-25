package http

import (
	v1 "antares-me/monitoring-system/internal/delivery/http/v1"
	"antares-me/monitoring-system/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	monitoringService service.MonitoringService
}

func NewHandler(monitoringService service.MonitoringService) *Handler {
	return &Handler{
		monitoringService: monitoringService,
	}
}

func (h *Handler) Init() *mux.Router {
	m := mux.NewRouter()

	h.InitMonitoring(m)

	return m
}

func (h *Handler) InitMonitoring(m *mux.Router) {
	handlerV1 := v1.NewHandler(h.monitoringService)
	handlerV1.Init(m)
}
