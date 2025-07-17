package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lcidral/goExpertOtel/pkg/models"
	"github.com/lcidral/goExpertOtel/pkg/telemetry"
	"github.com/lcidral/goExpertOtel/pkg/utils"
	"github.com/lcidral/goExpertOtel/services/service-a/internal/client"
)

// CEPHandler handles CEP-related HTTP requests
type CEPHandler struct {
	validator     *utils.CEPValidator
	serviceBClient *client.ServiceBClient
}

// NewCEPHandler creates a new CEP handler
func NewCEPHandler(serviceBClient *client.ServiceBClient) *CEPHandler {
	return &CEPHandler{
		validator:     utils.NewCEPValidator(),
		serviceBClient: serviceBClient,
	}
}

// HandleCEP processes CEP requests
func (h *CEPHandler) HandleCEP(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req models.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Erro ao decodificar requisição: %v", err)
		h.respondWithError(w, http.StatusBadRequest, models.ErrInvalidZipcode)
		return
	}

	// Start CEP validation span
	ctx, validationSpan := telemetry.StartSpan(r.Context(), "cep.validation",
		attribute.String("cep.input", req.CEP),
	)
	defer validationSpan.End()

	// Validate and normalize CEP
	normalizedCEP, err := h.validator.ValidateAndNormalize(req.CEP)
	if err != nil {
		log.Printf("CEP inválido: %s, erro: %v", req.CEP, err)
		validationSpan.RecordError(err)
		validationSpan.SetStatus(codes.Error, "CEP validation failed")
		h.respondWithError(w, http.StatusUnprocessableEntity, models.ErrInvalidZipcode)
		return
	}

	validationSpan.SetAttributes(
		attribute.String("cep.normalized", normalizedCEP),
		attribute.Bool("cep.valid", true),
	)
	validationSpan.SetStatus(codes.Ok, "CEP validation successful")

	// Create context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Start Service B call span
	ctx, serviceBSpan := telemetry.StartSpan(ctxWithTimeout, "service_b.call",
		attribute.String("cep.value", normalizedCEP),
		attribute.String("service.name", "service-b"),
	)
	defer serviceBSpan.End()

	// Call Service B
	tempResponse, err := h.serviceBClient.GetTemperature(ctx, normalizedCEP)
	if err != nil {
		log.Printf("Erro ao chamar Serviço B para CEP %s: %v", normalizedCEP, err)
		serviceBSpan.RecordError(err)
		serviceBSpan.SetStatus(codes.Error, "Service B call failed")
		
		// Check if it's a ServiceBError
		if serviceBErr, ok := err.(*client.ServiceBError); ok {
			serviceBSpan.SetAttributes(
				attribute.Int("http.status_code", serviceBErr.GetStatusCode()),
				attribute.String("error.message", serviceBErr.Message),
			)
			h.respondWithError(w, serviceBErr.GetStatusCode(), serviceBErr.Message)
			return
		}
		
		// Generic error
		h.respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}

	serviceBSpan.SetAttributes(
		attribute.String("city.name", tempResponse.City),
		attribute.Float64("temp_c", tempResponse.TempC),
		attribute.Float64("temp_f", tempResponse.TempF),
		attribute.Float64("temp_k", tempResponse.TempK),
	)
	serviceBSpan.SetStatus(codes.Ok, "Service B call successful")

	// Respond with success
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tempResponse); err != nil {
		log.Printf("Erro ao codificar resposta: %v", err)
	}

	log.Printf("CEP %s processado com sucesso", normalizedCEP)
}

// HealthCheck provides a health check endpoint
func (h *CEPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string{
		"status": "healthy",
		"service": "service-a",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	
	json.NewEncoder(w).Encode(response)
}

// respondWithError sends an error response
func (h *CEPHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorResp := models.ErrorResponse{
		Message: message,
	}
	json.NewEncoder(w).Encode(errorResp)
}