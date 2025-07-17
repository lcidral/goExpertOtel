# PRP-004: Docker e Docker Compose Setup

## Resumo
Containerizar todos os componentes do sistema e configurar um ambiente de desenvolvimento completo com Docker Compose, incluindo serviços, observabilidade e dependências.

## Motivação
A containerização garante consistência entre ambientes, facilita o desenvolvimento local e simplifica o deployment. Docker Compose permite orquestrar todos os serviços necessários com uma única comando, incluindo Zipkin e OTEL Collector.

## Descrição Detalhada

### Componentes a Containerizar

1. **Service A** - Input validation service
2. **Service B** - Temperature orchestration service  
3. **OTEL Collector** - Coleta e processamento de traces
4. **Zipkin** - Visualização de traces
5. **Networks** - Isolamento e comunicação entre serviços

### Estratégia de Build

- Multi-stage builds para otimização
- Cache de dependências Go
- Imagens mínimas com `scratch` ou `alpine`
- Health checks para todos os serviços
- Graceful shutdown

## Implementação Proposta

### Estrutura de Arquivos no Monorepo
```
goExpertOtel/                    # Raiz do monorepo
├── docker-compose.yml          # Compose principal
├── docker-compose.override.yml # Override para desenvolvimento
├── .env.example               # Template de variáveis
├── Makefile                   # Comandos úteis
├── services/
│   ├── service-a/
│   │   └── Dockerfile         # Dockerfile do Service A
│   └── service-b/
│       └── Dockerfile         # Dockerfile do Service B
├── deployments/
│   ├── otel-collector/
│   │   ├── Dockerfile
│   │   └── otel-collector-config.yaml
│   └── scripts/
│       ├── wait-for-it.sh
│       └── init-services.sh
└── build/
    └── docker/
        ├── base.dockerfile    # Imagem base compartilhada
        └── builder.dockerfile # Builder compartilhado
```

### Dockerfile Service A
```dockerfile
# Build stage - usando contexto do monorepo
FROM golang:1.21-alpine AS builder

WORKDIR /workspace

# Instalar dependências do sistema
RUN apk add --no-cache git

# Copiar arquivos do monorepo
COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY services/service-a/ ./services/service-a/

# Download dependencies
RUN go mod download

# Build binary com path correto do monorepo
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -o service-a \
    ./services/service-a/cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /workspace/service-a .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["./service-a"]
```

### Dockerfile Service B
```dockerfile
# Build stage - usando contexto do monorepo
FROM golang:1.21-alpine AS builder

WORKDIR /workspace

# Install git for private dependencies
RUN apk add --no-cache git

# Copiar arquivos do monorepo
COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY services/service-b/ ./services/service-b/

# Download dependencies
RUN go mod download

# Build com otimizações e path correto
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o service-b \
    ./services/service-b/cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# Copy binary
COPY --from=builder /workspace/service-b .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8081/health || exit 1

EXPOSE 8081

CMD ["./service-b"]
```

### Docker Compose Configuration
```yaml
version: '3.8'

services:
  service-a:
    build:
      context: .                    # Contexto do monorepo
      dockerfile: services/service-a/Dockerfile
    container_name: service-a
    ports:
      - "8080:8080"
    environment:
      - SERVICE_B_URL=http://service-b:8081
      - OTEL_SERVICE_NAME=service-a
      - OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
      - OTEL_EXPORTER_OTLP_INSECURE=true
    depends_on:
      service-b:
        condition: service_healthy
      otel-collector:
        condition: service_started
    networks:
      - app-network
    restart: unless-stopped

  service-b:
    build:
      context: .                    # Contexto do monorepo
      dockerfile: services/service-b/Dockerfile
    container_name: service-b
    ports:
      - "8081:8081"
    environment:
      - WEATHER_API_KEY=${WEATHER_API_KEY}
      - VIACEP_API_URL=https://viacep.com.br/ws
      - WEATHER_API_URL=http://api.weatherapi.com/v1
      - CACHE_TTL=3600s
      - OTEL_SERVICE_NAME=service-b
      - OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
      - OTEL_EXPORTER_OTLP_INSECURE=true
    depends_on:
      otel-collector:
        condition: service_started
    networks:
      - app-network
    restart: unless-stopped

  otel-collector:
    build:
      context: ./deployments/otel-collector
      dockerfile: Dockerfile
    container_name: otel-collector
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./deployments/otel-collector/otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP
      - "8888:8888"   # Prometheus metrics
    depends_on:
      - zipkin
    networks:
      - app-network
    restart: unless-stopped

  zipkin:
    image: openzipkin/zipkin:latest
    container_name: zipkin
    ports:
      - "9411:9411"
    environment:
      - STORAGE_TYPE=mem
      - MEM_MAX_SPANS=100000
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### Docker Compose Override (Development)
```yaml
version: '3.8'

services:
  service-a:
    build:
      context: .
      dockerfile: services/service-a/Dockerfile
      target: builder
    volumes:
      - .:/workspace                # Volume do monorepo completo
      - go-mod-cache:/go/pkg/mod   # Cache de módulos Go
    working_dir: /workspace
    command: go run ./services/service-a/cmd/server/main.go
    environment:
      - GO_ENV=development
      - CGO_ENABLED=0

  service-b:
    build:
      context: .
      dockerfile: services/service-b/Dockerfile
      target: builder
    volumes:
      - .:/workspace                # Volume do monorepo completo
      - go-mod-cache:/go/pkg/mod   # Cache de módulos Go
    working_dir: /workspace
    command: go run ./services/service-b/cmd/server/main.go
    environment:
      - GO_ENV=development
      - CGO_ENABLED=0

volumes:
  go-mod-cache:                     # Cache compartilhado de módulos Go
```

### Environment File (.env.example)
```env
# WeatherAPI Configuration
WEATHER_API_KEY=your_api_key_here

# Service Configuration
SERVICE_A_PORT=8080
SERVICE_B_PORT=8081

# OTEL Configuration
OTEL_COLLECTOR_PORT=4317
ZIPKIN_PORT=9411

# Development Settings
GO_ENV=development
LOG_LEVEL=debug
```

## Tarefas

### Setup Monorepo
- [ ] Configurar contextos de build para monorepo
- [ ] Criar Dockerfile para Service A com contexto correto
- [ ] Criar Dockerfile para Service B com contexto correto
- [ ] Configurar build cache compartilhado para pkg/
- [ ] Implementar health checks

### Docker Compose
- [ ] Criar docker-compose.yml com contextos do monorepo
- [ ] Criar docker-compose.override.yml para desenvolvimento
- [ ] Configurar volumes compartilhados para desenvolvimento
- [ ] Configurar networks e isolamento
- [ ] Definir ordem de inicialização
- [ ] Configurar volume para cache de módulos Go

### Infraestrutura
- [ ] Criar Dockerfile para OTEL Collector em deployments/
- [ ] Mover scripts para deployments/scripts/
- [ ] Configurar volumes para configurações
- [ ] Implementar init containers se necessário

### Scripts e Utilidades
- [ ] Criar script wait-for-it.sh em deployments/scripts/
- [ ] Criar Makefile na raiz do monorepo
- [ ] Implementar script de inicialização
- [ ] Criar script de backup/restore

### Configurações
- [ ] Criar .env.example na raiz do monorepo
- [ ] Documentar todas as variáveis
- [ ] Configurar logging centralizado
- [ ] Implementar secrets management

### Otimizações Monorepo
- [ ] Minimizar tamanho das imagens
- [ ] Implementar build cache para código compartilhado
- [ ] Configurar resource limits
- [ ] Adicionar labels para identificação do monorepo
- [ ] Otimizar copying de arquivos do monorepo

## Critérios de Aceitação

1. `docker-compose up` inicia todos os serviços
2. Todos os health checks passam
3. Serviços se comunicam corretamente
4. Traces aparecem no Zipkin
5. Hot reload funciona em desenvolvimento
6. Imagens otimizadas (< 50MB cada)
7. Graceful shutdown implementado
8. Documentação clara de uso

## Comandos Úteis (Makefile)

```makefile
# Development
up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

# Build
build:
	docker-compose build --no-cache

build-prod:
	docker-compose -f docker-compose.yml build

# Testing
test:
	docker-compose run --rm service-a go test ./...
	docker-compose run --rm service-b go test ./...

# Monitoring
zipkin:
	open http://localhost:9411

health:
	curl http://localhost:8080/health
	curl http://localhost:8081/health

# Cleanup
clean:
	docker-compose down -v
	docker system prune -f
```

## Riscos e Mitigações

| Risco | Mitigação |
|-------|-----------|
| Dependências não disponíveis | Health checks e restart policies |
| Conflitos de porta | Configuração customizável via .env |
| Performance em desenvolvimento | Volume mounts otimizados |
| Secrets em imagens | Multi-stage builds e runtime env |

## Boas Práticas Implementadas

1. **Segurança**
   - Non-root user nas imagens
   - Minimal base images
   - No secrets em build time
   - Network isolation

2. **Performance**
   - Build cache optimization
   - Layer caching
   - Minimal runtime dependencies
   - Resource limits

3. **Desenvolvimento**
   - Hot reload support
   - Debug configurations
   - Local volume mounts
   - Override configurations

## Monitoramento e Debug

### Logs Centralizados
```bash
# Ver logs de todos os serviços
docker-compose logs -f

# Logs específicos com timestamp
docker-compose logs -f --timestamps service-a

# Filtrar por nível
docker-compose logs -f | grep -E "ERROR|WARN"
```

### Métricas e Status
```bash
# Status dos containers
docker-compose ps

# Uso de recursos
docker stats

# Inspecionar rede
docker network inspect goexpertotel_app-network
```

## Referências
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Go Docker Multi-stage Builds](https://docs.docker.com/language/golang/build-images/)
- [Docker Compose Networking](https://docs.docker.com/compose/networking/)
- [Health Check Best Practices](https://docs.docker.com/engine/reference/builder/#healthcheck)