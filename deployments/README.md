# GoExpertOtel - Deployment

Este diretório contém as configurações de deployment para o sistema de temperatura por CEP com OpenTelemetry e Zipkin.

## Estrutura

```
deployments/
├── docker/                     # Docker Compose e Dockerfiles
│   ├── docker-compose.yml     # Orquestração completa
│   ├── Dockerfile.service-a   # Container Service A
│   ├── Dockerfile.service-b   # Container Service B
│   ├── .env.example          # Variáveis de ambiente exemplo
│   └── .env                  # Suas configurações (git ignored)
├── otel-collector/            # Configuração OpenTelemetry
│   └── otel-collector-config.yaml
└── README.md                 # Este arquivo
```

## Quick Start

### 1. Configurar API Key

```bash
cd deployments/docker
cp .env.example .env
# Edite .env e adicione sua WEATHER_API_KEY
```

### 2. Subir o Sistema

```bash
# A partir da raiz do projeto
docker-compose -f deployments/docker/docker-compose.yml up -d
```

### 3. Verificar Saúde

```bash
# Service A
curl http://localhost:8080/health

# Service B
curl http://localhost:8081/health

# Zipkin UI
open http://localhost:9411
```

### 4. Testar Funcionalidade

```bash
# Teste básico
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'
```

## Serviços

### Service A (Port 8080)
- **Função**: Entrada do sistema, validação de CEP
- **Endpoints**: 
  - `POST /` - Processar CEP
  - `GET /health` - Health check

### Service B (Port 8081)
- **Função**: Orquestração, APIs externas, conversão
- **Endpoints**:
  - `POST /temperature` - Obter temperatura
  - `GET /health` - Health check com stats de cache
  - `GET /cache/stats` - Estatísticas de cache

### OpenTelemetry Collector (Port 4317/4318)
- **Função**: Coleta e processa traces
- **Configuração**: `otel-collector/otel-collector-config.yaml`
- **Exporta para**: Zipkin e logs

### Zipkin (Port 9411)
- **Função**: Interface de visualização de traces
- **UI**: http://localhost:9411
- **Storage**: In-memory (desenvolvimento)

## Observabilidade

### Traces Implementados

#### Service A
- `http.request` - Requisição HTTP principal
- `cep.validation` - Validação de CEP
- `service_b.call` - Chamada para Service B

#### Service B
- `http.request` - Requisição HTTP principal
- `cep.validation` - Revalidação de CEP
- `cache.lookup` - Consultas ao cache
- `cache.store` - Armazenamento no cache
- `opencep.api.call` - Chamadas para OpenCEP
- `weather.api.call` - Chamadas para WeatherAPI
- `temperature.conversion` - Conversões de temperatura

### Visualizando Traces

1. **Acesse Zipkin**: http://localhost:9411
2. **Busque traces**: Use filtros como:
   - Por serviço: `service-a` ou `service-b`
   - Por tag: `cep.value=01310100`
   - Por erro: `error=true`
   - Por duração: `minDuration=100ms`

### Queries Úteis no Zipkin

```
# Traces com erro
error=true

# Traces lentos (> 1 segundo)
minDuration=1000ms

# Por CEP específico
tag=cep.value:01310100

# Falhas na WeatherAPI
serviceName=service-b AND tag=api.name:weather AND error=true

# Cache misses
tag=cache.hit:false

# Conversões de temperatura
operationName=temperature.conversion
```

## Monitoramento

### Logs Estruturados

```bash
# Logs dos serviços
docker-compose -f deployments/docker/docker-compose.yml logs -f service-a service-b

# Logs do Collector
docker-compose -f deployments/docker/docker-compose.yml logs -f otel-collector
```

### Health Checks

```bash
# Automatizado via Docker
docker-compose -f deployments/docker/docker-compose.yml ps

# Manual
curl http://localhost:8080/health
curl http://localhost:8081/health
```

### Métricas de Cache

```bash
# Estatísticas detalhadas
curl http://localhost:8081/cache/stats

# Incluído no health check
curl http://localhost:8081/health | jq .cache_stats
```

## Desenvolvimento

### Build Local

```bash
# Service A
go build -o bin/service-a ./services/service-a/cmd/server

# Service B  
go build -o bin/service-b ./services/service-b/cmd/server
```

### Desenvolvimento sem Docker

```bash
# Terminal 1 - Zipkin
docker run -d -p 9411:9411 openzipkin/zipkin:2.24

# Terminal 2 - OTEL Collector
docker run -p 4317:4317 -p 4318:4318 \
  -v $(pwd)/deployments/otel-collector/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
  otel/opentelemetry-collector-contrib:0.88.0 \
  --config=/etc/otel-collector-config.yaml

# Terminal 3 - Service B
export WEATHER_API_KEY=your_key_here
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
./bin/service-b

# Terminal 4 - Service A
export SERVICE_B_URL=http://localhost:8081
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
./bin/service-a
```

## Troubleshooting

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
   
   # Verifique conectividade
   curl http://localhost:4317/v1/traces -I
   ```

3. **Serviços não conseguem se comunicar**
   ```bash
   # Verifique network
   docker network ls
   docker-compose ps
   ```

### Reset Completo

```bash
# Parar e remover tudo
docker-compose -f deployments/docker/docker-compose.yml down -v

# Rebuild e restart
docker-compose -f deployments/docker/docker-compose.yml up --build -d
```

## Performance

### Configurações de Produção

Para produção, considere:

1. **Zipkin com Storage Persistente**
2. **OTEL Collector com Batch Otimizado**
3. **Sampling Rate Apropriado**
4. **Métricas Adicionais**
5. **Alertas de Observabilidade**

### Exemplo de Produção

```yaml
# docker-compose.prod.yml
services:
  zipkin:
    environment:
      - STORAGE_TYPE=elasticsearch
      - ES_HOSTS=http://elasticsearch:9200
    
  otel-collector:
    environment:
      - OTEL_TRACES_SAMPLER=probabilistic
      - OTEL_TRACES_SAMPLER_ARG=0.1  # 10% sampling
```

## Referências

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Zipkin Documentation](https://zipkin.io/pages/documentation.html)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Go OpenTelemetry](https://opentelemetry.io/docs/languages/go/)