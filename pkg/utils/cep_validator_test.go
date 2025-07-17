package utils

import (
	"testing"

	"github.com/lcidral/goExpertOtel/pkg/models"
)

func TestCEPValidator_ValidateCEP(t *testing.T) {
	validator := NewCEPValidator()

	tests := []struct {
		name    string
		cep     string
		wantErr bool
	}{
		{"CEP válido com 8 dígitos", "12345678", false},
		{"CEP válido - formato brasileiro comum", "01310100", false},
		{"CEP inválido - menos de 8 dígitos", "1234567", true},
		{"CEP inválido - mais de 8 dígitos", "123456789", true},
		{"CEP inválido - contém letras", "1234567a", true},
		{"CEP inválido - contém caracteres especiais", "12345-67", true},
		{"CEP inválido - string vazia", "", true},
		{"CEP inválido - apenas espaços", "        ", true},
		{"CEP válido - com espaços nas extremidades", "  12345678  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCEP(tt.cep)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateCEP() esperava erro, mas não recebeu nenhum")
				}
				if err != nil && err.Error() != models.ErrInvalidZipcode {
					t.Errorf("ValidateCEP() erro = %v, esperava %v", err.Error(), models.ErrInvalidZipcode)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCEP() erro inesperado = %v", err)
				}
			}
		})
	}
}

func TestCEPValidator_NormalizeCEP(t *testing.T) {
	validator := NewCEPValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"CEP sem formatação", "12345678", "12345678"},
		{"CEP com hífen", "12345-678", "12345678"},
		{"CEP com espaços", "123 456 78", "12345678"},
		{"CEP com pontos", "123.456.78", "12345678"},
		{"CEP com múltiplos caracteres especiais", "123-45.6 78", "12345678"},
		{"CEP com espaços nas extremidades", "  12345678  ", "12345678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.NormalizeCEP(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeCEP() = %v, esperava %v", result, tt.expected)
			}
		})
	}
}

func TestCEPValidator_ValidateAndNormalize(t *testing.T) {
	validator := NewCEPValidator()

	tests := []struct {
		name          string
		input         string
		expectedCEP   string
		expectedError bool
	}{
		{"CEP válido com hífen", "12345-678", "12345678", false},
		{"CEP válido com espaços", " 123 456 78 ", "12345678", false},
		{"CEP inválido após normalização - muito curto", "123-45", "", true},
		{"CEP inválido com letras", "123a5-678", "", true},
		{"CEP formato brasileiro válido", "01310-100", "01310100", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidateAndNormalize(tt.input)
			
			if tt.expectedError {
				if err == nil {
					t.Errorf("ValidateAndNormalize() esperava erro, mas não recebeu nenhum")
				}
				if result != "" {
					t.Errorf("ValidateAndNormalize() retornou CEP %v quando deveria estar vazio", result)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateAndNormalize() erro inesperado = %v", err)
				}
				if result != tt.expectedCEP {
					t.Errorf("ValidateAndNormalize() = %v, esperava %v", result, tt.expectedCEP)
				}
			}
		})
	}
}

// Benchmark para medir performance da validação
func BenchmarkCEPValidator_ValidateCEP(b *testing.B) {
	validator := NewCEPValidator()
	cep := "12345678"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateCEP(cep)
	}
}

func BenchmarkCEPValidator_NormalizeCEP(b *testing.B) {
	validator := NewCEPValidator()
	cep := "123-45.6 78"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.NormalizeCEP(cep)
	}
}