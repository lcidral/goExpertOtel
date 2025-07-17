# Sistema de Temperatura por CEP com OpenTelemetry

> 🌡️ Sistema distribuído em Go para consulta de temperatura por CEP brasileiro com observabilidade completa
> 
> ✨ Implementado como **monorepo** com OpenTelemetry e Zipkin para tracing distribuído

## 🏗️ Arquitetura

Este projeto implementa uma arquitetura de microserviços distribuída:

- **Service A**: Recebe e valida CEPs (porta 8080)
- **Service B**: Orquestra busca de localização e temperatura (porta 8081)
- **OpenTelemetry Collector**: Coleta e processa traces (porta 4317/4318)
- **Zipkin**: Interface de visualização de traces (porta 9411)

## 📁 Estrutura do Monorepo

```
goExpertOtel/
├── services/              # Microserviços
│   ├── service-a/         # Input e validação de CEP
│   │   ├── cmd/server/    # Aplicação principal
│   │   ├── internal/      # Lógica interna
│   │   │   ├── handler/   # HTTP handlers
│   │   │   └── client/    # Cliente HTTP para Service B
│   │   └── config/        # Configurações
│   └── service-b/         # Orquestração de temperatura
│       ├── cmd/server/    # Aplicação principal
│       ├── internal/      # Lógica interna
│       │   ├── handler/   # HTTP handlers
│       │   ├── client/    # Clientes APIs externas
│       │   ├── cache/     # Cache em memória
│       │   └── service/   # Conversões de temperatura
│       └── config/        # Configurações
├── pkg/                   # Código compartilhado
│   ├── models/           # Modelos de dados
│   ├── telemetry/        # OpenTelemetry compartilhado
│   └── utils/            # Utilitários (validação CEP)
├── deployments/          # Docker e infraestrutura
│   ├── docker/           # Docker Compose e Dockerfiles
│   └── otel-collector/   # Configuração OTEL Collector
└── PRPs/                 # Pull Request Proposals (documentação técnica)
```

## 🚀 Quick Start

### Pré-requisitos
- Docker e Docker Compose
- Go 1.21+ (para desenvolvimento local)
- Chave API do WeatherAPI (gratuita em [weatherapi.com](https://www.weatherapi.com/))

### Execução com Docker

```bash
# 1. Clone o repositório
git clone <repository-url>
cd goExpertOtel

# 2. Configure a API key
cd deployments/docker
cp .env.example .env
# Edite .env e adicione sua WEATHER_API_KEY

# 3. Inicie todo o sistema
docker-compose up -d

# 4. Verifique se os serviços estão funcionando
curl http://localhost:8080/health  # Service A
curl http://localhost:8081/health  # Service B
```

### Desenvolvimento Local

```bash
# Instalar dependências
go mod download

# Terminal 1 - Service B
cd services/service-b
export WEATHER_API_KEY=your_key_here
go run cmd/server/main.go

# Terminal 2 - Service A (em outro terminal)
cd services/service-a
export SERVICE_B_URL=http://localhost:8081
go run cmd/server/main.go
```

## 📡 Endpoints da API

### Service A (Entrada do Sistema)
- **Porta**: 8080
- **POST /**: Recebe CEP para consulta de temperatura
- **GET /health**: Health check

### Service B (Orquestração - Interno)
- **Porta**: 8081
- **POST /temperature**: Busca temperatura (chamado pelo Service A)
- **GET /health**: Health check com estatísticas de cache
- **GET /cache/stats**: Estatísticas detalhadas do cache

### Observabilidade
- **Zipkin UI**: http://localhost:9411 - Visualização de traces
- **OTEL Collector**: Porta 4317 (gRPC) / 4318 (HTTP)

## 🌡️ Exemplo de Uso

### CEP Válido
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Resposta (200 OK)
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

### CEP Inválido
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'

# Resposta (422 Unprocessable Entity)
{
  "message": "invalid zipcode"
}
```

### CEP Não Encontrado
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "00000000"}'

# Resposta (404 Not Found)
{
  "message": "can not find zipcode"
}
```

## 🔍 Observabilidade e Traces

### Spans Implementados

#### Service A:
- `http.request` - Requisição HTTP principal
- `cep.validation` - Validação de CEP
- `service_b.call` - Chamada para Service B

#### Service B:
- `http.request` - Requisição HTTP principal
- `cep.validation` - Revalidação de CEP
- `cache.lookup` / `cache.store` - Operações de cache
- `opencep.api.call` - Chamadas para OpenCEP API
- `weather.api.call` - Chamadas para WeatherAPI
- `temperature.conversion` - Conversões matemáticas

### Visualização no Zipkin

1. **Acesse**: http://localhost:9411
2. **Busque traces** usando:
   - Por serviço: `service-a` ou `service-b`
   - Por CEP: `cep.value=01310100`
   - Por erro: `error=true`
   - Por duração: `minDuration=100ms`
   - Cache misses: `cache.hit=false`

## 🧪 Testes

```bash
# Testes unitários (cobertura atual: ~85%)
go test ./... -v -cover

# Testes específicos
go test ./pkg/utils -v              # Validação de CEP
go test ./services/service-b/... -v # Conversões de temperatura

# Relatório de cobertura
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🐳 Docker e Infraestrutura

### Serviços no Docker Compose
- `service-a`: Service A containerizado
- `service-b`: Service B containerizado
- `otel-collector`: OpenTelemetry Collector
- `zipkin`: Interface de traces

### Comandos Úteis
```bash
# Ver logs
docker-compose -f deployments/docker/docker-compose.yml logs -f

# Parar sistema
docker-compose -f deployments/docker/docker-compose.yml down

# Rebuild
docker-compose -f deployments/docker/docker-compose.yml up --build
```

## ⚙️ Configuração

### Variáveis de Ambiente
```env
# Service A
PORT=8080
SERVICE_B_URL=http://service-b:8081
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317

# Service B
PORT=8081
WEATHER_API_KEY=your_api_key_here
OPENCEP_API_URL=https://opencep.com
WEATHER_API_URL=http://api.weatherapi.com/v1
CACHE_TTL=1h
CACHE_CLEANUP=10m

# OpenTelemetry
OTEL_SERVICE_NAME=service-a|service-b
OTEL_EXPORTER_OTLP_INSECURE=true
```

## 🚀 Features Implementadas

### ✅ Requisitos Atendidos
- [x] Service A: Validação de CEP e encaminhamento
- [x] Service B: Busca de localização e temperatura
- [x] Códigos HTTP corretos (200, 404, 422)
- [x] OpenTelemetry com tracing distribuído
- [x] Zipkin para visualização
- [x] Docker Compose para ambiente completo

### ✨ Features Extras
- [x] **Monorepo** com código compartilhado
- [x] **Cache inteligente** com TTL diferenciado
- [x] **Health checks** com estatísticas
- [x] **Graceful shutdown**
- [x] **Logs estruturados**
- [x] **Testes unitários** abrangentes
- [x] **Instrumentation HTTP** automática
- [x] **Spans customizados** detalhados

## 🎯 Cache Strategy

O Service B implementa cache em múltiplas camadas:

1. **Cache de Localização** (24h TTL): Dados do OpenCEP
2. **Cache de Temperatura** (10min TTL): Dados meteorológicos
3. **Cache de Resposta** (10min TTL): Resposta processada completa

## 🔧 Troubleshooting

### Problemas Comuns

1. **WEATHER_API_KEY não configurada**
   ```bash
   # Verifique se está no .env
   grep WEATHER_API_KEY deployments/docker/.env
   ```

2. **Traces não aparecem no Zipkin**
   ```bash
   # Verifique logs do collector
   docker-compose logs otel-collector
   ```

3. **Service B falha ao buscar clima**
   - Verifique se a API key do WeatherAPI está válida
   - Confirme conectividade com APIs externas

### Health Checks
```bash
# Verificar saúde dos serviços
curl http://localhost:8080/health
curl http://localhost:8081/health

# Estatísticas de cache
curl http://localhost:8081/cache/stats
```

## 📚 Documentação Técnica

- [Documentação de Deployment](deployments/README.md)
- [Service A README](services/service-a/README.md)
- [Service B README](services/service-b/README.md)
- [PRPs (Pull Request Proposals)](PRPs/) - Documentação técnica detalhada

## 🏆 Especificações Técnicas

### APIs Externas Utilizadas
- **OpenCEP**: `https://opencep.com/v1/{cep}.json` - Busca de localização
- **WeatherAPI**: `http://api.weatherapi.com/v1/current.json` - Dados meteorológicos

### Conversões de Temperatura
- **Celsius → Fahrenheit**: `F = C × 1.8 + 32`
- **Celsius → Kelvin**: `K = C + 273`
- **Precisão**: 1 casa decimal

### Performance
- **Cache hit rate**: Monitore via `/cache/stats`
- **Response time**: Visível nos traces do Zipkin
- **Concurrent requests**: Suporte nativo com Go goroutines

## 🤝 Contribuição

O projeto está estruturado como monorepo para facilitar desenvolvimento e manutenção:

1. Código compartilhado em `pkg/`
2. Serviços independentes em `services/`
3. Infraestrutura centralizada em `deployments/`
4. Documentação técnica em `PRPs/`

---

**Desenvolvido com ❤️ em Go** • **OpenTelemetry** • **Docker** • **Zipkin**