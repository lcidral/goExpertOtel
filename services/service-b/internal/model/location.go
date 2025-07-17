package model

// ViaCEPResponse representa a resposta da API de CEP (compatível com OpenCEP e ViaCEP)
type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade,omitempty"`    // OpenCEP only
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Estado      string `json:"estado,omitempty"`     // OpenCEP only
	Regiao      string `json:"regiao,omitempty"`     // OpenCEP only
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia,omitempty"`        // ViaCEP only
	DDD         string `json:"ddd,omitempty"`        // ViaCEP only
	SIAFI       string `json:"siafi,omitempty"`      // ViaCEP only
	Erro        bool   `json:"erro,omitempty"`       // ViaCEP only
}

// IsValid verifica se a resposta da API de CEP é válida
func (v *ViaCEPResponse) IsValid() bool {
	return !v.Erro && v.Localidade != "" && v.UF != ""
}

// GetFullLocation retorna a localização completa no formato "Cidade, Estado"
func (v *ViaCEPResponse) GetFullLocation() string {
	if v.UF != "" {
		return v.Localidade + ", " + v.UF
	}
	return v.Localidade
}

// GetCityName retorna apenas o nome da cidade
func (v *ViaCEPResponse) GetCityName() string {
	return v.Localidade
}