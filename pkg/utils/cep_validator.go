package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lcidral/goExpertOtel/pkg/models"
)

// CEPValidator implementa validação de CEP brasileiro
type CEPValidator struct {
	cepRegex *regexp.Regexp
}

// NewCEPValidator cria uma nova instância do validador
func NewCEPValidator() *CEPValidator {
	// Regex para validar CEP: exatamente 8 dígitos
	cepRegex := regexp.MustCompile(`^\d{8}$`)
	
	return &CEPValidator{
		cepRegex: cepRegex,
	}
}

// ValidateCEP valida se o CEP está no formato correto
func (v *CEPValidator) ValidateCEP(cep string) error {
	if cep == "" {
		return fmt.Errorf(models.ErrInvalidZipcode)
	}

	// Remove espaços em branco
	cep = strings.TrimSpace(cep)
	
	// Valida o formato: deve ter exatamente 8 dígitos
	if !v.cepRegex.MatchString(cep) {
		return fmt.Errorf(models.ErrInvalidZipcode)
	}

	return nil
}

// NormalizeCEP normaliza o CEP removendo caracteres especiais
func (v *CEPValidator) NormalizeCEP(cep string) string {
	// Remove espaços, hífens e outros caracteres não numéricos
	normalized := regexp.MustCompile(`[^0-9]`).ReplaceAllString(cep, "")
	return normalized
}

// ValidateAndNormalize valida e normaliza o CEP em uma única operação
func (v *CEPValidator) ValidateAndNormalize(cep string) (string, error) {
	// Primeiro normaliza
	normalized := v.NormalizeCEP(cep)
	
	// Depois valida
	if err := v.ValidateCEP(normalized); err != nil {
		return "", err
	}
	
	return normalized, nil
}