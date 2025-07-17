package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/lcidral/goExpertOtel/pkg/models"
)

// ServiceBClient cliente para comunicação com o Serviço B
type ServiceBClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewServiceBClient cria uma nova instância do cliente
func NewServiceBClient(baseURL string, timeout time.Duration) *ServiceBClient {
	return &ServiceBClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// GetTemperature faz uma requisição para o Serviço B para obter temperatura
func (c *ServiceBClient) GetTemperature(ctx context.Context, cep string) (*models.TemperatureResponse, error) {
	// Prepara o payload
	request := models.CEPRequest{
		CEP: cep,
	}

	// Serializa para JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Cria a requisição HTTP
	url := fmt.Sprintf("%s/temperature", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Define headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Executa a requisição
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Verifica status HTTP
	switch resp.StatusCode {
	case http.StatusOK:
		// Sucesso - decodifica resposta
		var tempResponse models.TemperatureResponse
		if err := json.NewDecoder(resp.Body).Decode(&tempResponse); err != nil {
			return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
		}
		return &tempResponse, nil

	case http.StatusNotFound:
		// CEP não encontrado
		return nil, &ServiceBError{
			StatusCode: resp.StatusCode,
			Message:    models.ErrZipcodeNotFound,
		}

	case http.StatusUnprocessableEntity:
		// CEP inválido
		return nil, &ServiceBError{
			StatusCode: resp.StatusCode,
			Message:    models.ErrInvalidZipcode,
		}

	default:
		// Outros erros
		var errorResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, &ServiceBError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("erro inesperado do serviço B: %d", resp.StatusCode),
			}
		}
		return nil, &ServiceBError{
			StatusCode: resp.StatusCode,
			Message:    errorResp.Message,
		}
	}
}

// ServiceBError representa erros específicos do Serviço B
type ServiceBError struct {
	StatusCode int
	Message    string
}

func (e *ServiceBError) Error() string {
	return e.Message
}

func (e *ServiceBError) GetStatusCode() int {
	return e.StatusCode
}