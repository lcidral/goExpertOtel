package models

// CEPRequest representa a estrutura da requisição com CEP
type CEPRequest struct {
	CEP string `json:"cep" validate:"required,len=8,numeric"`
}

// TemperatureResponse representa a estrutura da resposta com informações de temperatura
type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ErrorResponse representa a estrutura de resposta de erro
type ErrorResponse struct {
	Message string `json:"message"`
}

// Constants para mensagens de erro padrão
const (
	ErrInvalidZipcode  = "invalid zipcode"
	ErrZipcodeNotFound = "can not find zipcode"
)
