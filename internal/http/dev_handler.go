package httpapi

import (
	"net/http"

	devapp "github.com/mohit838/olario-platform-backend/internal/application/dev"
)

// DevHandler exposes local-only learning and smoke-test endpoints.
// These endpoints are not product APIs; they exist to prove the full request
// path through HTTP, application service, Postgres, Redis, and JSON response.
type DevHandler struct {
	service *devapp.Service
}

func NewDevHandler(service *devapp.Service) *DevHandler {
	return &DevHandler{service: service}
}

// FullCircle seeds one complete grocery flow and returns an invoice-style
// response. It is intentionally POST because it writes demo data.
func (h *DevHandler) FullCircle(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "dev service is unavailable; check Postgres and Redis config",
		})
		return
	}

	result, err := h.service.RunFullCircle(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, result)
}
