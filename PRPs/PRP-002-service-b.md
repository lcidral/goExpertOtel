# PRP-002: Implementação do Serviço B - Temperature Orchestration Service

## Resumo
Implementar o serviço de orquestração responsável por buscar informações de localização através do CEP e obter dados de temperatura, retornando as temperaturas em Celsius, Fahrenheit e Kelvin.

## Motivação
O Serviço B é o núcleo do sistema, responsável por orquestrar as chamadas às APIs externas (ViaCEP e WeatherAPI) e realizar as conversões de temperatura necessárias. Sua implementação correta é crucial para a funcionalidade principal do sistema.

## Descrição Detalhada

### Arquitetura
- Servidor HTTP em Go
- Porta padrão: 8081
- Integração com APIs externas:
  - ViaCEP para busca de localização
  - WeatherAPI para dados meteorológicos

### Funcionalidades

1. **Recebimento de CEP**
   - Validação do formato (8 dígitos)
   - Processamento assíncrono quando possível

2. **Busca de Localização (ViaCEP)**
   - Endpoint: `https://viacep.com.br/ws/{cep}/json/`
   - Cache de resultados para otimização
   - Tratamento de CEPs não encontrados

3. **Busca de Temperatura (WeatherAPI)**
   - Endpoint: `http://api.weatherapi.com/v1/current.json`
   - Parâmetros: API key e localização
   - Extração da temperatura em Celsius

4. **Conversões de Temperatura**
   - Celsius → Fahrenheit: `F = C * 1.8 + 32`
   - Celsius → Kelvin: `K = C + 273`
   - Precisão: 1 casa decimal

5. **Respostas HTTP**
   - 200: Sucesso com dados completos
   - 404: CEP não encontrado ("can not find zipcode")
   - 422: CEP inválido ("invalid zipcode")

## Implementação Proposta

### Estrutura de Monorepo
O Serviço B faz parte do monorepo e compartilha código comum com outros serviços.

```
goExpertOtel/                    # Raiz do monorepo
├── services/
│   ├── service-a/              # Serviço A (implementado no PRP-001)
│   └── service-b/              # Serviço B - Temperature Service
│       ├── cmd/
│       │   └── server/
│       │       └── main.go
│       ├── internal/
│       │   ├── handler/
│       │   │   └── temperature_handler.go
│       │   ├── service/
│       │   │   ├── location_service.go
│       │   │   ├── weather_service.go
│       │   │   └── temperature_converter.go
│       │   ├── client/
│       │   │   ├── viacep_client.go
│       │   │   └── weather_client.go
│       │   ├── cache/
│       │   │   └── memory_cache.go
│       │   └── model/
│       │       ├── location.go
│       │       └── temperature.go
│       └── config/
│           └── config.go
├── pkg/                        # Código compartilhado entre serviços
│   ├── models/
│   │   ├── cep.go             # Modelo CEP compartilhado
│   │   └── temperature.go     # Modelo Temperature compartilhado
│   ├── telemetry/             # Telemetria compartilhada (PRP-003)
│   └── utils/
│       └── http.go            # Utilitários HTTP compartilhados
├── go.mod                     # Módulo Go único para o monorepo
├── go.sum
└── docker-compose.yml
```

### Module Path e Imports
Utilizando o módulo principal do monorepo:

```go
// Imports do Service B
import (
    "github.com/yourusername/goExpertOtel/pkg/models"
    "github.com/yourusername/goExpertOtel/pkg/telemetry"
    "github.com/yourusername/goExpertOtel/pkg/utils"
)
```

### Dependências
- `net/http` - Cliente e servidor HTTP
- `encoding/json` - Parsing JSON
- `github.com/go-chi/chi/v5` - Router HTTP
- `github.com/patrickmn/go-cache` - Cache em memória
- `github.com/joho/godotenv` - Variáveis de ambiente

### Configurações
```env
PORT=8081
WEATHER_API_KEY=your_api_key_here
WEATHER_API_URL=http://api.weatherapi.com/v1
VIACEP_API_URL=https://viacep.com.br/ws
CACHE_TTL=3600s
REQUEST_TIMEOUT=10s
```

### Modelos de Dados

```go
// Modelos compartilhados em pkg/models/

// pkg/models/cep.go
type CEPRequest struct {
    CEP string `json:"cep"`
}

// pkg/models/temperature.go
type TemperatureResponse struct {
    City  string  `json:"city"`
    TempC float64 `json:"temp_C"`
    TempF float64 `json:"temp_F"`
    TempK float64 `json:"temp_K"`
}

// Modelos internos do Service B

// internal/model/location.go
type ViaCEPResponse struct {
    CEP         string `json:"cep"`
    Localidade  string `json:"localidade"`
    UF          string `json:"uf"`
    Erro        bool   `json:"erro"`
}

// internal/model/weather.go
type WeatherAPIResponse struct {
    Current struct {
        TempC float64 `json:"temp_c"`
    } `json:"current"`
}
```

## Tarefas

- [ ] Criar estrutura services/service-b no monorepo
- [ ] Implementar modelos compartilhados em pkg/models
- [ ] Criar modelos internos específicos do serviço
- [ ] Criar cliente para ViaCEP API
- [ ] Criar cliente para WeatherAPI
- [ ] Implementar serviço de conversão de temperaturas
- [ ] Criar cache em memória com TTL configurável
- [ ] Implementar handler principal de temperatura
- [ ] Adicionar middleware de logging e recovery
- [ ] Criar testes unitários para conversões
- [ ] Criar testes de integração com mocks das APIs
- [ ] Implementar circuit breaker para APIs externas
- [ ] Adicionar métricas de performance
- [ ] Implementar rate limiting para APIs externas
- [ ] Criar health check com verificação de APIs
- [ ] Configurar Dockerfile com build context do monorepo

## Critérios de Aceitação

1. Busca correta de localização via CEP
2. Obtenção precisa de temperatura atual
3. Conversões matemáticas corretas (verificadas com testes)
4. Respostas HTTP adequadas para cada cenário
5. Cache funcionando corretamente
6. Tratamento de erros das APIs externas
7. Logs estruturados para debugging
8. Testes com cobertura mínima de 85%

## Riscos e Mitigações

| Risco | Mitigação |
|-------|-----------|
| Limite de rate das APIs externas | Implementar cache agressivo e rate limiting |
| APIs externas indisponíveis | Circuit breaker e respostas em cache |
| CEPs válidos mas sem dados meteorológicos | Fallback para cidade mais próxima |
| Custos da WeatherAPI | Monitorar uso e implementar quotas |
| Dados desatualizados no cache | TTL configurável e invalidação manual |

## Considerações de Performance

1. **Cache Strategy**
   - Cache de CEP → Localização (TTL: 24h)
   - Cache de Localização → Temperatura (TTL: 10min)
   - Cache warming para CEPs populares

2. **Otimizações**
   - Connection pooling para HTTP clients
   - Timeouts agressivos
   - Processamento paralelo quando possível

## Referências
- [ViaCEP Documentation](https://viacep.com.br/)
- [WeatherAPI Documentation](https://www.weatherapi.com/docs/)
- [Go Caching Best Practices](https://github.com/patrickmn/go-cache)
- [Circuit Breaker Pattern](https://github.com/sony/gobreaker)