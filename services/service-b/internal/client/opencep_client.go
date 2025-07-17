package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lcidral/goExpertOtel/services/service-b/internal/model"
)

// OpenCEPClient cliente para a API OpenCEP
type OpenCEPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewOpenCEPClient cria uma nova instância do cliente OpenCEP
func NewOpenCEPClient(baseURL string, timeout time.Duration) *OpenCEPClient {
	return &OpenCEPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetLocationByCEP busca informações de localização pelo CEP
func (c *OpenCEPClient) GetLocationByCEP(ctx context.Context, cep string) (*model.ViaCEPResponse, error) {
	// Constrói a URL da API (OpenCEP usa formato /v1/{cep}.json)
	url := fmt.Sprintf("%s/v1/%s.json", c.baseURL, cep)

	// Cria a requisição HTTP
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, &CEPNotFoundError{CEP: cep}
		}
		return nil, fmt.Errorf("erro na API OpenCEP: status %d", resp.StatusCode)
	}

	// Decodifica a resposta
	var openCepResp model.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&openCepResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// OpenCEP não tem campo "erro", verifica se campos essenciais estão presentes
	if openCepResp.Localidade == "" || openCepResp.UF == "" {
		return nil, &CEPNotFoundError{CEP: cep}
	}

	// Valida se a resposta contém dados essenciais
	if !openCepResp.IsValid() {
		return nil, fmt.Errorf("resposta inválida da API OpenCEP para CEP %s", cep)
	}

	return &openCepResp, nil
}

// CEPNotFoundError representa erro quando CEP não é encontrado
type CEPNotFoundError struct {
	CEP string
}

func (e *CEPNotFoundError) Error() string {
	return fmt.Sprintf("CEP %s não encontrado", e.CEP)
}

func (e *CEPNotFoundError) IsCEPNotFound() bool {
	return true
}