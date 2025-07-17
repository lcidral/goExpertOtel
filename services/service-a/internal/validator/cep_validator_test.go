package validator

import (
	"testing"

	"github.com/lcidral/goExpertOtel/pkg/models"
	"github.com/lcidral/goExpertOtel/pkg/utils"
)

func TestCEPValidator_ValidateCEP(t *testing.T) {
	validator := utils.NewCEPValidator()

	tests := []struct {
		name    string
		cep     string
		wantErr bool
	}{
		{
			name:    "CEP válido com 8 dígitos",
			cep:     "12345678",
			wantErr: false,
		},
		{
			name:    "CEP válido - formato brasileiro comum",
			cep:     "01310100",
			wantErr: false,
		},
		{
			name:    "CEP inválido - menos de 8 dígitos",
			cep:     "1234567",
			wantErr: true,
		},
		{
			name:    "CEP inválido - mais de 8 dígitos",
			cep:     "123456789",
			wantErr: true,
		},
		{
			name:    "CEP inválido - contém letras",
			cep:     "1234567a",
			wantErr: true,
		},
		{
			name:    "CEP inválido - contém caracteres especiais",
			cep:     "12345-67",
			wantErr: true,
		},
		{
			name:    "CEP inválido - string vazia",
			cep:     "",
			wantErr: true,
		},
		{
			name:    "CEP inválido - apenas espaços",
			cep:     "        ",
			wantErr: true,
		},
		{
			name:    "CEP válido - com espaços nas extremidades",
			cep:     "  12345678  ",
			wantErr: false,
		},
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
	validator := utils.NewCEPValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CEP sem formatação",
			input:    "12345678",
			expected: "12345678",
		},
		{
			name:     "CEP com hífen",
			input:    "12345-678",
			expected: "12345678",
		},
		{
			name:     "CEP com espaços",
			input:    "123 456 78",
			expected: "12345678",
		},
		{
			name:     "CEP com pontos",
			input:    "123.456.78",
			expected: "12345678",
		},
		{
			name:     "CEP com múltiplos caracteres especiais",
			input:    "123-45.6 78",
			expected: "12345678",
		},
		{
			name:     "CEP com espaços nas extremidades",
			input:    "  12345678  ",
			expected: "12345678",
		},
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
	validator := utils.NewCEPValidator()

	tests := []struct {
		name          string
		input         string
		expectedCEP   string
		expectedError bool
	}{
		{
			name:          "CEP válido com hífen",
			input:         "12345-678",
			expectedCEP:   "12345678",
			expectedError: false,
		},
		{
			name:          "CEP válido com espaços",
			input:         " 123 456 78 ",
			expectedCEP:   "12345678",
			expectedError: false,
		},
		{
			name:          "CEP inválido após normalização - muito curto",
			input:         "123-45",
			expectedCEP:   "",
			expectedError: true,
		},
		{
			name:          "CEP inválido com letras",
			input:         "123a5-678",
			expectedCEP:   "",
			expectedError: true,
		},
		{
			name:          "CEP formato brasileiro válido",
			input:         "01310-100",
			expectedCEP:   "01310100",
			expectedError: false,
		},
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
	validator := utils.NewCEPValidator()
	cep := "12345678"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateCEP(cep)
	}
}

func BenchmarkCEPValidator_NormalizeCEP(b *testing.B) {
	validator := utils.NewCEPValidator()
	cep := "123-45.6 78"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.NormalizeCEP(cep)
	}
}