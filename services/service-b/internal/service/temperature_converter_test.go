package service

import (
	"testing"
)

func TestTemperatureConverter_CelsiusToFahrenheit(t *testing.T) {
	converter := NewTemperatureConverter()

	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{"Zero Celsius", 0, 32},
		{"Ponto de ebulição da água", 100, 212},
		{"Temperatura ambiente", 25, 77},
		{"Temperatura negativa", -10, 14},
		{"Temperatura corporal", 36.5, 97.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.CelsiusToFahrenheit(tt.celsius)
			if result != tt.expected {
				t.Errorf("CelsiusToFahrenheit(%v) = %v, esperava %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestTemperatureConverter_CelsiusToKelvin(t *testing.T) {
	converter := NewTemperatureConverter()

	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{"Zero absoluto", -273, 0},
		{"Zero Celsius", 0, 273},
		{"Ponto de ebulição da água", 100, 373},
		{"Temperatura ambiente", 25, 298},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.CelsiusToKelvin(tt.celsius)
			if result != tt.expected {
				t.Errorf("CelsiusToKelvin(%v) = %v, esperava %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestTemperatureConverter_FahrenheitToCelsius(t *testing.T) {
	converter := NewTemperatureConverter()

	tests := []struct {
		name       string
		fahrenheit float64
		expected   float64
	}{
		{"Congelamento da água", 32, 0},
		{"Ebulição da água", 212, 100},
		{"Temperatura ambiente", 77, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.FahrenheitToCelsius(tt.fahrenheit)
			// Usar tolerância para comparação de float
			if abs(result-tt.expected) > 0.01 {
				t.Errorf("FahrenheitToCelsius(%v) = %v, esperava %v", tt.fahrenheit, result, tt.expected)
			}
		})
	}
}

func TestTemperatureConverter_RoundToOneDecimal(t *testing.T) {
	converter := NewTemperatureConverter()

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"Sem arredondamento necessário", 25.5, 25.5},
		{"Arredondar para cima", 25.56, 25.6},
		{"Arredondar para baixo", 25.54, 25.5},
		{"Número inteiro", 25, 25.0},
		{"Múltiplas casas decimais", 25.123456, 25.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.roundToOneDecimal(tt.input)
			if result != tt.expected {
				t.Errorf("roundToOneDecimal(%v) = %v, esperava %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTemperatureConverter_ConvertToAllUnits(t *testing.T) {
	converter := NewTemperatureConverter()

	tests := []struct {
		name     string
		celsius  float64
		city     string
		tempC    float64
		tempF    float64
		tempK    float64
	}{
		{
			name:    "Temperatura ambiente em São Paulo",
			celsius: 25.0,
			city:    "São Paulo",
			tempC:   25.0,
			tempF:   77.0,
			tempK:   298.0,
		},
		{
			name:    "Temperatura fria no Rio de Janeiro",
			celsius: 10.5,
			city:    "Rio de Janeiro",
			tempC:   10.5,
			tempF:   50.9,
			tempK:   283.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertToAllUnits(tt.celsius, tt.city)

			if result.City != tt.city {
				t.Errorf("City = %v, esperava %v", result.City, tt.city)
			}
			if result.TempC != tt.tempC {
				t.Errorf("TempC = %v, esperava %v", result.TempC, tt.tempC)
			}
			if result.TempF != tt.tempF {
				t.Errorf("TempF = %v, esperava %v", result.TempF, tt.tempF)
			}
			if result.TempK != tt.tempK {
				t.Errorf("TempK = %v, esperava %v", result.TempK, tt.tempK)
			}
		})
	}
}

func TestTemperatureConverter_ValidateTemperature(t *testing.T) {
	converter := NewTemperatureConverter()

	tests := []struct {
		name     string
		celsius  float64
		expected bool
	}{
		{"Temperatura válida normal", 25, true},
		{"Temperatura válida fria", -10, true},
		{"Temperatura válida quente", 50, true},
		{"Temperatura válida limite inferior", -100, true},
		{"Temperatura válida limite superior", 60, true},
		{"Temperatura inválida muito fria", -101, false},
		{"Temperatura inválida muito quente", 61, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ValidateTemperature(tt.celsius)
			if result != tt.expected {
				t.Errorf("ValidateTemperature(%v) = %v, esperava %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestTemperatureConverter_GetTemperatureDescription(t *testing.T) {
	converter := NewTemperatureConverter()

	tests := []struct {
		name     string
		celsius  float64
		expected string
	}{
		{"Muito frio", -5, "muito frio"},
		{"Frio", 5, "frio"},
		{"Fresco", 15, "fresco"},
		{"Agradável", 22, "agradável"},
		{"Quente", 28, "quente"},
		{"Muito quente", 32, "muito quente"},
		{"Extremamente quente", 40, "extremamente quente"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.GetTemperatureDescription(tt.celsius)
			if result != tt.expected {
				t.Errorf("GetTemperatureDescription(%v) = %v, esperava %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

// abs retorna o valor absoluto de um float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Benchmark tests
func BenchmarkTemperatureConverter_CelsiusToFahrenheit(b *testing.B) {
	converter := NewTemperatureConverter()
	celsius := 25.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.CelsiusToFahrenheit(celsius)
	}
}

func BenchmarkTemperatureConverter_ConvertToAllUnits(b *testing.B) {
	converter := NewTemperatureConverter()
	celsius := 25.0
	city := "São Paulo"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.ConvertToAllUnits(celsius, city)
	}
}