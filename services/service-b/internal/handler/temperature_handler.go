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
	"github.com/lcidral/goExpertOtel/services/service-b/internal/cache"
	"github.com/lcidral/goExpertOtel/services/service-b/internal/client"
	"github.com/lcidral/goExpertOtel/services/service-b/internal/model"
	"github.com/lcidral/goExpertOtel/services/service-b/internal/service"
)

// TemperatureHandler handles temperature-related HTTP requests
type TemperatureHandler struct {
	openCEPClient    *client.OpenCEPClient
	weatherClient    *client.WeatherClient
	tempConverter    *service.TemperatureConverter
	cache            *cache.MemoryCache
	validator        *utils.CEPValidator
}

// NewTemperatureHandler creates a new temperature handler
func NewTemperatureHandler(
	openCEPClient *client.OpenCEPClient,
	weatherClient *client.WeatherClient,
	cache *cache.MemoryCache,
) *TemperatureHandler {
	return &TemperatureHandler{
		openCEPClient: openCEPClient,
		weatherClient: weatherClient,
		tempConverter: service.NewTemperatureConverter(),
		cache:         cache,
		validator:     utils.NewCEPValidator(),
	}
}

// HandleTemperature processes temperature requests
func (h *TemperatureHandler) HandleTemperature(w http.ResponseWriter, r *http.Request) {
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

	// Check cache first for complete temperature response
	ctx, cacheSpan := telemetry.StartSpan(ctx, "cache.lookup",
		attribute.String("cache.key", "temp:"+normalizedCEP),
		attribute.String("cache.type", "temperature"),
	)
	if cachedTemp, found := h.cache.GetTemperature(normalizedCEP); found {
		log.Printf("Cache hit para temperatura do CEP %s", normalizedCEP)
		cacheSpan.SetAttributes(
			attribute.Bool("cache.hit", true),
			attribute.String("city.name", cachedTemp.City),
		)
		cacheSpan.SetStatus(codes.Ok, "Cache hit")
		cacheSpan.End()
		h.respondWithSuccess(w, cachedTemp)
		return
	}
	cacheSpan.SetAttributes(attribute.Bool("cache.hit", false))
	cacheSpan.SetStatus(codes.Ok, "Cache miss")
	cacheSpan.End()

	// Create context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get location from CEP
	location, err := h.getLocationWithCache(ctxWithTimeout, normalizedCEP)
	if err != nil {
		log.Printf("Erro ao buscar localização para CEP %s: %v", normalizedCEP, err)
		
		// Check if it's a CEP not found error
		if _, isCEPNotFound := err.(*client.CEPNotFoundError); isCEPNotFound {
			h.respondWithError(w, http.StatusNotFound, models.ErrZipcodeNotFound)
			return
		}
		
		h.respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}

	// Get weather data
	weather, err := h.getWeatherWithCache(ctxWithTimeout, location.GetFullLocation())
	if err != nil {
		log.Printf("Erro ao buscar clima para %s: %v", location.GetFullLocation(), err)
		
		// Check if it's a location not found error
		if _, isLocationNotFound := err.(*client.LocationNotFoundError); isLocationNotFound {
			h.respondWithError(w, http.StatusNotFound, models.ErrZipcodeNotFound)
			return
		}
		
		h.respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}

	// Start temperature conversion span
	_, conversionSpan := telemetry.StartSpan(ctxWithTimeout, "temperature.conversion",
		attribute.Float64("temp_c", weather.GetTemperatureCelsius()),
		attribute.String("city.name", location.GetCityName()),
	)
	defer conversionSpan.End()

	// Convert temperatures
	tempResponse := h.tempConverter.ConvertToAllUnits(
		weather.GetTemperatureCelsius(),
		location.GetCityName(),
	)

	conversionSpan.SetAttributes(
		attribute.Float64("temp_f", tempResponse.TempF),
		attribute.Float64("temp_k", tempResponse.TempK),
	)
	conversionSpan.SetStatus(codes.Ok, "Temperature conversion successful")

	// Cache the complete response
	h.cache.SetTemperature(normalizedCEP, tempResponse, 10*time.Minute)

	// Respond with success
	h.respondWithSuccess(w, tempResponse)
	log.Printf("CEP %s processado com sucesso: %s, %.1f°C", 
		normalizedCEP, tempResponse.City, tempResponse.TempC)
}

// getLocationWithCache gets location with caching
func (h *TemperatureHandler) getLocationWithCache(ctx context.Context, cep string) (*model.ViaCEPResponse, error) {
	// Start OpenCEP API call span
	ctx, opencepSpan := telemetry.StartSpan(ctx, "opencep.api.call",
		attribute.String("cep.value", cep),
		attribute.String("api.name", "opencep"),
	)
	defer opencepSpan.End()

	// Check cache first
	ctx, cacheSpan := telemetry.StartSpan(ctx, "cache.lookup",
		attribute.String("cache.key", "location:"+cep),
		attribute.String("cache.type", "location"),
	)
	if cachedLocation, found := h.cache.GetLocation(cep); found {
		log.Printf("Cache hit para localização do CEP %s", cep)
		cacheSpan.SetAttributes(
			attribute.Bool("cache.hit", true),
			attribute.String("city.name", cachedLocation.GetCityName()),
		)
		cacheSpan.SetStatus(codes.Ok, "Cache hit")
		cacheSpan.End()
		
		opencepSpan.SetAttributes(
			attribute.Bool("cache.hit", true),
			attribute.String("city.name", cachedLocation.GetCityName()),
		)
		opencepSpan.SetStatus(codes.Ok, "Location retrieved from cache")
		return cachedLocation, nil
	}
	cacheSpan.SetAttributes(attribute.Bool("cache.hit", false))
	cacheSpan.SetStatus(codes.Ok, "Cache miss")
	cacheSpan.End()

	// Not in cache, fetch from API
	location, err := h.openCEPClient.GetLocationByCEP(ctx, cep)
	if err != nil {
		opencepSpan.RecordError(err)
		opencepSpan.SetStatus(codes.Error, "OpenCEP API call failed")
		return nil, err
	}

	// Cache the result (cache for 24 hours for location data)
	_, cacheStoreSpan := telemetry.StartSpan(ctx, "cache.store",
		attribute.String("cache.key", "location:"+cep),
		attribute.String("cache.type", "location"),
		attribute.String("cache.ttl", "24h"),
	)
	h.cache.SetLocation(cep, location, 24*time.Hour)
	cacheStoreSpan.SetStatus(codes.Ok, "Location cached")
	cacheStoreSpan.End()

	opencepSpan.SetAttributes(
		attribute.Bool("cache.hit", false),
		attribute.String("city.name", location.GetCityName()),
		attribute.String("state", location.UF),
	)
	opencepSpan.SetStatus(codes.Ok, "Location retrieved from OpenCEP API")
	log.Printf("Cache miss para localização do CEP %s, dados armazenados", cep)

	return location, nil
}

// getWeatherWithCache gets weather with caching
func (h *TemperatureHandler) getWeatherWithCache(ctx context.Context, location string) (*model.WeatherAPIResponse, error) {
	// Start Weather API call span
	ctx, weatherSpan := telemetry.StartSpan(ctx, "weather.api.call",
		attribute.String("location", location),
		attribute.String("api.name", "weather"),
	)
	defer weatherSpan.End()

	// Check cache first
	ctx, cacheSpan := telemetry.StartSpan(ctx, "cache.lookup",
		attribute.String("cache.key", "weather:"+location),
		attribute.String("cache.type", "weather"),
	)
	if cachedWeather, found := h.cache.GetWeather(location); found {
		log.Printf("Cache hit para clima de %s", location)
		cacheSpan.SetAttributes(
			attribute.Bool("cache.hit", true),
			attribute.Float64("temp_c", cachedWeather.GetTemperatureCelsius()),
		)
		cacheSpan.SetStatus(codes.Ok, "Cache hit")
		cacheSpan.End()
		
		weatherSpan.SetAttributes(
			attribute.Bool("cache.hit", true),
			attribute.Float64("temp_c", cachedWeather.GetTemperatureCelsius()),
		)
		weatherSpan.SetStatus(codes.Ok, "Weather retrieved from cache")
		return cachedWeather, nil
	}
	cacheSpan.SetAttributes(attribute.Bool("cache.hit", false))
	cacheSpan.SetStatus(codes.Ok, "Cache miss")
	cacheSpan.End()

	// Not in cache, fetch from API
	weather, err := h.weatherClient.GetCurrentWeather(ctx, location)
	if err != nil {
		weatherSpan.RecordError(err)
		weatherSpan.SetStatus(codes.Error, "Weather API call failed")
		return nil, err
	}

	// Cache the result (cache for 10 minutes for weather data)
	_, cacheStoreSpan := telemetry.StartSpan(ctx, "cache.store",
		attribute.String("cache.key", "weather:"+location),
		attribute.String("cache.type", "weather"),
		attribute.String("cache.ttl", "10m"),
	)
	h.cache.SetWeather(location, weather, 10*time.Minute)
	cacheStoreSpan.SetStatus(codes.Ok, "Weather cached")
	cacheStoreSpan.End()

	weatherSpan.SetAttributes(
		attribute.Bool("cache.hit", false),
		attribute.Float64("temp_c", weather.GetTemperatureCelsius()),
		attribute.String("condition", weather.Current.Condition.Text),
	)
	weatherSpan.SetStatus(codes.Ok, "Weather retrieved from API")
	log.Printf("Cache miss para clima de %s, dados armazenados", location)

	return weather, nil
}

// HealthCheck provides a health check endpoint
func (h *TemperatureHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	cacheStats := h.cache.Stats()
	
	response := map[string]interface{}{
		"status":      "healthy",
		"service":     "service-b",
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"cache_stats": cacheStats,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CacheStats provides cache statistics endpoint
func (h *TemperatureHandler) CacheStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	stats := h.cache.Stats()
	json.NewEncoder(w).Encode(stats)
}

// respondWithSuccess sends a success response
func (h *TemperatureHandler) respondWithSuccess(w http.ResponseWriter, data *models.TemperatureResponse) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erro ao codificar resposta: %v", err)
	}
}

// respondWithError sends an error response
func (h *TemperatureHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorResp := models.ErrorResponse{
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("Erro ao codificar resposta de erro: %v", err)
	}
}