package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lcidral/goExpertOtel/services/service-b/internal/model"
)

// WeatherClient cliente para a WeatherAPI
type WeatherClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewWeatherClient cria uma nova instância do cliente WeatherAPI
func NewWeatherClient(baseURL, apiKey string, timeout time.Duration) *WeatherClient {
	return &WeatherClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetCurrentWeather busca informações meteorológicas atuais para uma localização
func (c *WeatherClient) GetCurrentWeather(ctx context.Context, location string) (*model.WeatherAPIResponse, error) {
	// Constrói a URL da API
	endpoint := fmt.Sprintf("%s/current.json", c.baseURL)
	
	// Cria os parâmetros da query
	params := url.Values{}
	params.Add("key", c.apiKey)
	params.Add("q", location)
	params.Add("aqi", "no") // Não precisamos de dados de qualidade do ar

	// URL completa
	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	// Cria a requisição HTTP
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Define headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "goExpertOtel-service-b/1.0")

	// Executa a requisição
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Verifica status HTTP
	switch resp.StatusCode {
	case http.StatusOK:
		// Sucesso - decodifica resposta normal
		var weatherResp model.WeatherAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
			return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
		}

		// Valida se a resposta contém dados essenciais
		if !weatherResp.IsValid() {
			return nil, fmt.Errorf("resposta inválida da WeatherAPI para localização %s", location)
		}

		return &weatherResp, nil

	case http.StatusBadRequest:
		// Erro 400 - localização não encontrada ou inválida
		var errorResp model.WeatherAPIError
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, &LocationNotFoundError{Location: location}
		}
		return nil, &LocationNotFoundError{
			Location: location,
			Message:  errorResp.GetMessage(),
		}

	case http.StatusUnauthorized:
		// Erro 401 - API key inválida
		return nil, fmt.Errorf("API key inválida para WeatherAPI")

	case http.StatusForbidden:
		// Erro 403 - quota excedida ou acesso negado
		return nil, fmt.Errorf("quota excedida ou acesso negado na WeatherAPI")

	default:
		// Outros erros
		return nil, fmt.Errorf("erro na WeatherAPI: status %d", resp.StatusCode)
	}
}

// GetCurrentWeatherByCoordinates busca informações meteorológicas por coordenadas
func (c *WeatherClient) GetCurrentWeatherByCoordinates(ctx context.Context, lat, lon float64) (*model.WeatherAPIResponse, error) {
	location := fmt.Sprintf("%.6f,%.6f", lat, lon)
	return c.GetCurrentWeather(ctx, location)
}

// LocationNotFoundError representa erro quando localização não é encontrada
type LocationNotFoundError struct {
	Location string
	Message  string
}

func (e *LocationNotFoundError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("localização '%s' não encontrada: %s", e.Location, e.Message)
	}
	return fmt.Sprintf("localização '%s' não encontrada", e.Location)
}

func (e *LocationNotFoundError) IsLocationNotFound() bool {
	return true
}