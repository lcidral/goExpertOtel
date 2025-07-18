package model

// WeatherAPIResponse representa a resposta completa da WeatherAPI
type WeatherAPIResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int64   `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		UV         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

// GetTemperatureCelsius retorna a temperatura em Celsius
func (w *WeatherAPIResponse) GetTemperatureCelsius() float64 {
	return w.Current.TempC
}

// GetLocationName retorna o nome da localização
func (w *WeatherAPIResponse) GetLocationName() string {
	return w.Location.Name
}

// IsValid verifica se a resposta da WeatherAPI é válida
func (w *WeatherAPIResponse) IsValid() bool {
	return w.Location.Name != "" && w.Current.TempC != 0
}

// WeatherAPIError representa erros específicos da WeatherAPI
type WeatherAPIError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *WeatherAPIError) GetMessage() string {
	return e.Error.Message
}

func (e *WeatherAPIError) GetCode() int {
	return e.Error.Code
}