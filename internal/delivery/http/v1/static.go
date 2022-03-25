package v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) StaticInit(m *mux.Router) {
	fs := http.FileServer(http.Dir("./web/"))
	m.PathPrefix("/").Handler(fs)
}
