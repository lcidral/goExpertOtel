package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/lcidral/goExpertOtel/pkg/models"
	"github.com/lcidral/goExpertOtel/services/service-b/internal/model"
)

// MemoryCache cache em memória para otimizar chamadas às APIs
type MemoryCache struct {
	cache *cache.Cache
}

// NewMemoryCache cria uma nova instância do cache
func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) *MemoryCache {
	return &MemoryCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Cache Keys patterns
const (
	locationCacheKey = "location:%s"    // location:12345678
	weatherCacheKey  = "weather:%s"     // weather:São Paulo,SP
	tempCacheKey     = "temp:%s"        // temp:12345678
)

// GetLocation busca localização no cache
func (mc *MemoryCache) GetLocation(cep string) (*model.ViaCEPResponse, bool) {
	key := fmt.Sprintf(locationCacheKey, cep)
	if item, found := mc.cache.Get(key); found {
		if location, ok := item.(*model.ViaCEPResponse); ok {
			return location, true
		}
	}
	return nil, false
}

// SetLocation armazena localização no cache
func (mc *MemoryCache) SetLocation(cep string, location *model.ViaCEPResponse, duration time.Duration) {
	key := fmt.Sprintf(locationCacheKey, cep)
	mc.cache.Set(key, location, duration)
}

// GetWeather busca dados meteorológicos no cache
func (mc *MemoryCache) GetWeather(location string) (*model.WeatherAPIResponse, bool) {
	key := fmt.Sprintf(weatherCacheKey, location)
	if item, found := mc.cache.Get(key); found {
		if weather, ok := item.(*model.WeatherAPIResponse); ok {
			return weather, true
		}
	}
	return nil, false
}

// SetWeather armazena dados meteorológicos no cache
func (mc *MemoryCache) SetWeather(location string, weather *model.WeatherAPIResponse, duration time.Duration) {
	key := fmt.Sprintf(weatherCacheKey, location)
	mc.cache.Set(key, weather, duration)
}

// GetTemperature busca temperatura completa no cache
func (mc *MemoryCache) GetTemperature(cep string) (*models.TemperatureResponse, bool) {
	key := fmt.Sprintf(tempCacheKey, cep)
	if item, found := mc.cache.Get(key); found {
		if temp, ok := item.(*models.TemperatureResponse); ok {
			return temp, true
		}
	}
	return nil, false
}

// SetTemperature armazena temperatura completa no cache
func (mc *MemoryCache) SetTemperature(cep string, temp *models.TemperatureResponse, duration time.Duration) {
	key := fmt.Sprintf(tempCacheKey, cep)
	mc.cache.Set(key, temp, duration)
}

// InvalidateLocation remove localização do cache
func (mc *MemoryCache) InvalidateLocation(cep string) {
	key := fmt.Sprintf(locationCacheKey, cep)
	mc.cache.Delete(key)
}

// InvalidateWeather remove dados meteorológicos do cache
func (mc *MemoryCache) InvalidateWeather(location string) {
	key := fmt.Sprintf(weatherCacheKey, location)
	mc.cache.Delete(key)
}

// InvalidateTemperature remove temperatura do cache
func (mc *MemoryCache) InvalidateTemperature(cep string) {
	key := fmt.Sprintf(tempCacheKey, cep)
	mc.cache.Delete(key)
}

// Clear limpa todo o cache
func (mc *MemoryCache) Clear() {
	mc.cache.Flush()
}

// Stats retorna estatísticas do cache
func (mc *MemoryCache) Stats() map[string]interface{} {
	items := mc.cache.Items()
	locationCount := 0
	weatherCount := 0
	tempCount := 0

	for key := range items {
		switch {
		case len(key) > 9 && key[:9] == "location:":
			locationCount++
		case len(key) > 8 && key[:8] == "weather:":
			weatherCount++
		case len(key) > 5 && key[:5] == "temp:":
			tempCount++
		}
	}

	return map[string]interface{}{
		"total_items":    len(items),
		"location_items": locationCount,
		"weather_items":  weatherCount,
		"temp_items":     tempCount,
	}
}

// GetSize retorna o número total de itens no cache
func (mc *MemoryCache) GetSize() int {
	return mc.cache.ItemCount()
}