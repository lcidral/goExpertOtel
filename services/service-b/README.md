# Service B - Temperature Orchestration Service

Serviço de orquestração responsável por buscar informações de localização através do CEP e obter dados de temperatura, retornando as temperaturas em Celsius, Fahrenheit e Kelvin.

## Funcionalidades

- **Busca de Localização**: Integração com OpenCEP para obter dados de localização por CEP
- **Consulta de Temperatura**: Integração com WeatherAPI para obter dados meteorológicos
- **Conversão de Temperatura**: Converte automaticamente para Celsius, Fahrenheit e Kelvin
- **Cache Inteligente**: Cache em memória para otimizar chamadas às APIs externas
- **Health Check**: Endpoint de monitoramento com estatísticas de cache
- **Validação**: Revalidação de CEPs recebidos

## Endpoints

### POST /temperature
Recebe CEP e retorna informações de temperatura completas.

**Request:**
```json
{
  "cep": "01310100"
}
```

**Success Response (200):**
```json
{
  "city": "São Paulo",
  "temp_C": 25.5,
  "temp_F": 77.9,
  "temp_K": 298.5
}
```

**Error Responses:**
- `422`: CEP inválido (`{"message": "invalid zipcode"}`)
- `404`: CEP não encontrado (`{"message": "can not find zipcode"}`)
- `500`: Erro interno (APIs externas indisponíveis)

### GET /health
Endpoint de health check com estatísticas de cache.

**Response (200):**
```json
{
  "status": "healthy",
  "service": "service-b",
  "timestamp": "2024-01-01T12:00:00Z",
  "cache_stats": {
    "total_items": 15,
    "location_items": 5,
    "weather_items": 8,
    "temp_items": 2
  }
}
```

### GET /cache/stats
Estatísticas detalhadas do cache.

**Response (200):**
```json
{
  "total_items": 15,
  "location_items": 5,
  "weather_items": 8,
  "temp_items": 2
}
```

## Configuração

### Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PORT` | `8081` | Porta do serviço |
| `WEATHER_API_KEY` | **obrigatória** | Chave da WeatherAPI |
| `WEATHER_API_URL` | `http://api.weatherapi.com/v1` | URL base da WeatherAPI |
| `OPENCEP_API_URL` | `https://opencep.com` | URL base da OpenCEP |
| `REQUEST_TIMEOUT` | `10s` | Timeout para APIs externas |
| `CACHE_TTL` | `1h` | TTL padrão do cache |
| `CACHE_CLEANUP` | `10m` | Intervalo de limpeza do cache |

### APIs Externas

#### OpenCEP
- **URL**: `https://opencep.com/v1/{cep}.json`
- **Cache**: 24 horas (localização não muda frequentemente)
- **Gratuita**: Sem necessidade de API key

#### WeatherAPI
- **URL**: `http://api.weatherapi.com/v1/current.json`
- **Cache**: 10 minutos (dados meteorológicos mudam rapidamente)
- **Requer**: API key gratuita em [weatherapi.com](https://www.weatherapi.com/)

## Execução

### Desenvolvimento Local
```bash
# Configure a API key
export WEATHER_API_KEY="sua_chave_aqui"

# Da raiz do monorepo
cd services/service-b
go run cmd/server/main.go
```

### Build
```bash
# Da raiz do monorepo
go build -o bin/service-b ./services/service-b/cmd/server
./bin/service-b
```

## Cache Strategy

O serviço implementa cache em múltiplas camadas:

1. **Cache de Localização** (24h TTL):
   - Key: `location:{cep}`
   - Valor: Dados do OpenCEP
   - Justificativa: Localização não muda

2. **Cache de Temperatura** (10min TTL):
   - Key: `weather:{cidade,estado}`
   - Valor: Dados meteorológicos
   - Justificativa: Dados mudam rapidamente

3. **Cache de Resposta Completa** (10min TTL):
   - Key: `temp:{cep}`
   - Valor: Resposta final processada
   - Justificativa: Evita reprocessamento

## Conversões de Temperatura

Implementa conversões matemáticas precisas:

- **Celsius → Fahrenheit**: `F = C × 1.8 + 32`
- **Celsius → Kelvin**: `K = C + 273`
- **Precisão**: 1 casa decimal

## Testes

```bash
# Testes unitários
go test ./services/service-b/... -v

# Testes específicos do conversor
go test ./services/service-b/internal/service -v

# Testes com cobertura
go test ./services/service-b/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Monitoramento

### Logs Estruturados
```bash
# Service B iniciado
2024/01/01 12:00:00 🌡️ Service B iniciado na porta 8081

# Cache hits/misses
2024/01/01 12:00:01 Cache hit para localização do CEP 01310100
2024/01/01 12:00:02 Cache miss para clima de São Paulo,SP, dados armazenados

# Processamento
2024/01/01 12:00:03 CEP 01310100 processado com sucesso: São Paulo, 25.5°C
```

### Métricas de Cache
- **Hit Rate**: Monitore via `/cache/stats`
- **Memory Usage**: Cache em memória com limpeza automática
- **TTL Strategy**: Otimizado para diferentes tipos de dados

## Exemplos de Uso

```bash
# CEP válido
curl -X POST http://localhost:8081/temperature \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Health check
curl http://localhost:8081/health

# Estatísticas de cache
curl http://localhost:8081/cache/stats
```

## Tratamento de Erros

### Cenários Cobertos
1. **CEP Inválido**: Validação local antes das APIs
2. **CEP Não Encontrado**: OpenCEP retorna erro
3. **Localização Não Encontrada**: WeatherAPI não encontra cidade
4. **APIs Indisponíveis**: Timeout e erro de rede
5. **API Key Inválida**: WeatherAPI authentication error

### Resiliência
- **Timeouts**: 10s por requisição às APIs
- **Cache**: Reduz dependência das APIs externas
- **Graceful Degradation**: Logs detalhados para debugging

## Arquitetura

```
services/service-b/
├── cmd/server/              # Aplicação principal
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── service/            # Lógica de negócio
│   ├── client/             # Clientes das APIs
│   ├── cache/              # Cache em memória
│   └── model/              # Modelos internos
└── config/                 # Configurações
```

## Dependências

- `github.com/go-chi/chi/v5` - Router HTTP
- `github.com/patrickmn/go-cache` - Cache em memória
- `github.com/joho/godotenv` - Variáveis de ambiente
- Compartilha `pkg/models` e `pkg/utils` com outros serviços