package models

import (
	"encoding/json"
	"testing"
)

func TestCEPRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		request  CEPRequest
		expected string
	}{
		{
			name:     "CEP válido",
			request:  CEPRequest{CEP: "12345678"},
			expected: `{"cep":"12345678"}`,
		},
		{
			name:     "CEP vazio",
			request:  CEPRequest{CEP: ""},
			expected: `{"cep":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("Erro ao serializar JSON: %v", err)
			}

			if string(jsonData) != tt.expected {
				t.Errorf("JSON = %v, esperava %v", string(jsonData), tt.expected)
			}
		})
	}
}

func TestCEPRequest_JSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected CEPRequest
		wantErr  bool
	}{
		{
			name:     "JSON válido",
			jsonStr:  `{"cep":"12345678"}`,
			expected: CEPRequest{CEP: "12345678"},
			wantErr:  false,
		},
		{
			name:     "JSON com CEP vazio",
			jsonStr:  `{"cep":""}`,
			expected: CEPRequest{CEP: ""},
			wantErr:  false,
		},
		{
			name:     "JSON inválido",
			jsonStr:  `{"cep":}`,
			expected: CEPRequest{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request CEPRequest
			err := json.Unmarshal([]byte(tt.jsonStr), &request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Esperava erro, mas não recebeu nenhum")
				}
			} else {
				if err != nil {
					t.Errorf("Erro inesperado: %v", err)
				}
				if request.CEP != tt.expected.CEP {
					t.Errorf("CEP = %v, esperava %v", request.CEP, tt.expected.CEP)
				}
			}
		})
	}
}

func TestTemperatureResponse_JSONMarshaling(t *testing.T) {
	response := TemperatureResponse{
		City:  "São Paulo",
		TempC: 25.5,
		TempF: 77.9,
		TempK: 298.65,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Erro ao serializar JSON: %v", err)
	}

	expected := `{"city":"São Paulo","temp_C":25.5,"temp_F":77.9,"temp_K":298.65}`
	if string(jsonData) != expected {
		t.Errorf("JSON = %v, esperava %v", string(jsonData), expected)
	}
}

func TestTemperatureResponse_JSONUnmarshaling(t *testing.T) {
	jsonStr := `{"city":"São Paulo","temp_C":25.5,"temp_F":77.9,"temp_K":298.65}`

	var response TemperatureResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Errorf("Erro ao deserializar JSON: %v", err)
	}

	if response.City != "São Paulo" {
		t.Errorf("City = %v, esperava 'São Paulo'", response.City)
	}
	if response.TempC != 25.5 {
		t.Errorf("TempC = %v, esperava 25.5", response.TempC)
	}
	if response.TempF != 77.9 {
		t.Errorf("TempF = %v, esperava 77.9", response.TempF)
	}
	if response.TempK != 298.65 {
		t.Errorf("TempK = %v, esperava 298.65", response.TempK)
	}
}

func TestErrorResponse_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		response ErrorResponse
		expected string
	}{
		{
			name:     "Erro de CEP inválido",
			response: ErrorResponse{Message: ErrInvalidZipcode},
			expected: `{"message":"invalid zipcode"}`,
		},
		{
			name:     "Erro de CEP não encontrado",
			response: ErrorResponse{Message: ErrZipcodeNotFound},
			expected: `{"message":"can not find zipcode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.response)
			if err != nil {
				t.Errorf("Erro ao serializar JSON: %v", err)
			}

			if string(jsonData) != tt.expected {
				t.Errorf("JSON = %v, esperava %v", string(jsonData), tt.expected)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if ErrInvalidZipcode != "invalid zipcode" {
		t.Errorf("ErrInvalidZipcode = %v, esperava 'invalid zipcode'", ErrInvalidZipcode)
	}

	if ErrZipcodeNotFound != "can not find zipcode" {
		t.Errorf("ErrZipcodeNotFound = %v, esperava 'can not find zipcode'", ErrZipcodeNotFound)
	}
}
