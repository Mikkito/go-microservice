package handlers

import (
	"net/http"

	"go-microservice/internal/metrics"
)

type IntegrationHandler struct{}

func NewIntegrationHandler() *IntegrationHandler {
	return &IntegrationHandler{}
}

func (h *IntegrationHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	metrics.Handler().ServeHTTP(w, r)
}
