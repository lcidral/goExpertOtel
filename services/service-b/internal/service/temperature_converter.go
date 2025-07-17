package service

import (
	"math"

	"github.com/lcidral/goExpertOtel/pkg/models"
)

// TemperatureConverter serviço para conversão de temperaturas
type TemperatureConverter struct{}

// NewTemperatureConverter cria uma nova instância do conversor
func NewTemperatureConverter() *TemperatureConverter {
	return &TemperatureConverter{}
}

// ConvertToAllUnits converte temperatura de Celsius para todas as unidades
func (tc *TemperatureConverter) ConvertToAllUnits(celsius float64, city string) *models.TemperatureResponse {
	return &models.TemperatureResponse{
		City:  city,
		TempC: tc.roundToOneDecimal(celsius),
		TempF: tc.roundToOneDecimal(tc.CelsiusToFahrenheit(celsius)),
		TempK: tc.roundToOneDecimal(tc.CelsiusToKelvin(celsius)),
	}
}

// CelsiusToFahrenheit converte Celsius para Fahrenheit
// Fórmula: F = C * 1.8 + 32
func (tc *TemperatureConverter) CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// CelsiusToKelvin converte Celsius para Kelvin
// Fórmula: K = C + 273.15 (usando 273 conforme especificação)
func (tc *TemperatureConverter) CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

// FahrenheitToCelsius converte Fahrenheit para Celsius
// Fórmula: C = (F - 32) / 1.8
func (tc *TemperatureConverter) FahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) / 1.8
}

// KelvinToCelsius converte Kelvin para Celsius
// Fórmula: C = K - 273.15 (usando 273 conforme especificação)
func (tc *TemperatureConverter) KelvinToCelsius(kelvin float64) float64 {
	return kelvin - 273
}

// roundToOneDecimal arredonda para uma casa decimal
func (tc *TemperatureConverter) roundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

// ValidateTemperature valida se a temperatura está dentro de limites razoáveis
func (tc *TemperatureConverter) ValidateTemperature(celsius float64) bool {
	// Limites razoáveis para temperatura terrestre: -100°C a 60°C
	return celsius >= -100 && celsius <= 60
}

// GetTemperatureDescription retorna uma descrição textual da temperatura
func (tc *TemperatureConverter) GetTemperatureDescription(celsius float64) string {
	switch {
	case celsius < 0:
		return "muito frio"
	case celsius < 10:
		return "frio"
	case celsius < 20:
		return "fresco"
	case celsius < 25:
		return "agradável"
	case celsius < 30:
		return "quente"
	case celsius < 35:
		return "muito quente"
	default:
		return "extremamente quente"
	}
}