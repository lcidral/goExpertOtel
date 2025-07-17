# PRP-005: Documentação e Testes

## Resumo
Criar documentação abrangente do sistema, implementar suite completa de testes (unitários, integração e E2E) e estabelecer práticas de qualidade contínua.

## Motivação
Documentação clara e testes robustos são essenciais para a manutenibilidade e confiabilidade do sistema. Este PRP garante que o projeto seja facilmente compreensível, testável e extensível por qualquer desenvolvedor.

## Descrição Detalhada

### Componentes de Documentação

1. **README Principal**
   - Visão geral do projeto
   - Instruções de instalação e execução
   - Exemplos de uso
   - Troubleshooting comum

2. **Documentação Técnica**
   - Arquitetura do sistema
   - Decisões de design
   - Fluxos de dados
   - Diagramas de sequência

3. **API Documentation**
   - OpenAPI/Swagger specs
   - Exemplos de requests/responses
   - Códigos de erro
   - Rate limits

4. **Guias de Desenvolvimento**
   - Setup do ambiente
   - Convenções de código
   - Processo de contribuição
   - CI/CD pipeline

### Estratégia de Testes

1. **Testes Unitários**
   - Cobertura mínima: 80%
   - Foco em lógica de negócio
   - Mocks para dependências externas

2. **Testes de Integração**
   - Testes de API endpoints
   - Integração entre serviços
   - Testes com containers

3. **Testes E2E**
   - Fluxo completo do sistema
   - Validação com Zipkin
   - Performance benchmarks

## Implementação Proposta

### Estrutura de Documentação no Monorepo
```
goExpertOtel/                    # Raiz do monorepo
├── README.md                    # README principal do projeto
├── docs/
│   ├── ARCHITECTURE.md          # Arquitetura detalhada
│   ├── API.md                   # Documentação da API
│   ├── DEVELOPMENT.md           # Guia de desenvolvimento
│   ├── TROUBLESHOOTING.md       # Solução de problemas
│   ├── CONTRIBUTING.md          # Como contribuir
│   ├── MONOREPO.md             # Guia específico do monorepo
│   └── diagrams/
│       ├── architecture.puml    # Diagrama de arquitetura
│       ├── sequence.puml        # Diagramas de sequência
│       └── deployment.puml      # Diagrama de deployment
├── api/
│   ├── openapi.yaml            # Especificação OpenAPI
│   └── postman/
│       └── collection.json     # Coleção Postman
├── examples/
│   ├── requests/
│   │   ├── valid_cep.sh       # Exemplo de CEP válido
│   │   ├── invalid_cep.sh     # Exemplo de CEP inválido
│   │   └── not_found_cep.sh   # Exemplo de CEP não encontrado
│   └── responses/
│       ├── success.json
│       ├── invalid.json
│       └── not_found.json
├── services/
│   ├── service-a/
│   │   └── README.md           # Documentação específica do Service A
│   └── service-b/
│       └── README.md           # Documentação específica do Service B
└── pkg/
    └── README.md               # Documentação dos packages compartilhados
```

### Estrutura de Testes no Monorepo
```
goExpertOtel/                    # Raiz do monorepo
├── services/
│   ├── service-a/
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   │   └── cep_handler_test.go
│   │   │   ├── validator/
│   │   │   │   └── cep_validator_test.go
│   │   │   └── client/
│   │   │       └── service_b_client_test.go
│   │   └── test/
│   │       ├── integration/
│   │       │   └── api_test.go
│   │       └── e2e/
│   │           └── flow_test.go
│   └── service-b/
│       ├── internal/
│       │   ├── service/
│       │   │   ├── location_service_test.go
│       │   │   ├── weather_service_test.go
│       │   │   └── temperature_converter_test.go
│       │   └── client/
│       │       ├── viacep_client_test.go
│       │       └── weather_client_test.go
│       └── test/
│           ├── integration/
│           │   └── api_test.go
│           ├── mocks/
│           │   ├── viacep_mock.go
│           │   └── weather_mock.go
│           └── fixtures/
│               └── test_data.json
├── pkg/
│   ├── models/
│   │   └── models_test.go       # Testes dos modelos compartilhados
│   ├── telemetry/
│   │   └── telemetry_test.go    # Testes da telemetria compartilhada
│   └── utils/
│       └── utils_test.go        # Testes dos utilitários
├── test/
│   ├── e2e/                     # Testes E2E do sistema completo
│   │   ├── full_flow_test.go
│   │   └── zipkin_integration_test.go
│   ├── shared/                  # Código de teste compartilhado
│   │   ├── testcontainers.go
│   │   └── fixtures.go
│   └── load/                    # Testes de carga
│       └── performance_test.go
└── scripts/
    └── test/
        ├── run-all-tests.sh     # Script para rodar todos os testes
        └── coverage-report.sh   # Script para gerar relatório de cobertura
```

### README.md Principal (Monorepo)
```markdown
# Sistema de Temperatura por CEP com OpenTelemetry

> 🌡️ Sistema distribuído em Go para consulta de temperatura por CEP brasileiro
> 
> ✨ Implementado como **monorepo** com observabilidade completa usando OpenTelemetry e Zipkin

## 🏗️ Arquitetura

Este projeto implementa uma arquitetura de microserviços distribuída:

- **Service A**: Recebe e valida CEPs
- **Service B**: Orquestra busca de localização e temperatura
- **OpenTelemetry**: Observabilidade distribuída
- **Zipkin**: Visualização de traces

## 📁 Estrutura do Monorepo

```
goExpertOtel/
├── services/          # Microserviços
│   ├── service-a/     # Input e validação de CEP
│   └── service-b/     # Orquestração de temperatura
├── pkg/               # Código compartilhado
│   ├── models/        # Modelos de dados
│   ├── telemetry/     # OpenTelemetry
│   └── utils/         # Utilitários
├── deployments/       # Docker e infraestrutura
├── docs/              # Documentação
└── test/              # Testes E2E e compartilhados
```

## 🚀 Quick Start

### Pré-requisitos
- Docker e Docker Compose
- Go 1.21+ (para desenvolvimento)
- Chave API do WeatherAPI

### Instalação
\`\`\`bash
# 1. Clone o monorepo
git clone <repository-url>
cd goExpertOtel

# 2. Configure ambiente
cp .env.example .env
# Edite .env com sua WEATHER_API_KEY

# 3. Inicie todos os serviços
docker-compose up
\`\`\`

### Desenvolvimento Local
\`\`\`bash
# Instalar dependências do monorepo
go mod download

# Executar Service A
cd services/service-a && go run cmd/server/main.go

# Executar Service B (em outro terminal)
cd services/service-b && go run cmd/server/main.go
\`\`\`

## 📡 Endpoints

| Serviço | Porta | Endpoint | Descrição |
|---------|-------|----------|-----------|
| Service A | 8080 | POST / | Recebe CEP para consulta |
| Service B | 8081 | POST /temperature | Busca temperatura (interno) |
| Zipkin | 9411 | - | Interface de traces |

### Exemplo de Uso
\`\`\`bash
# CEP válido
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Resposta
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
\`\`\`

## 🔍 Observabilidade

- **Zipkin UI**: http://localhost:9411
- **OTEL Collector**: http://localhost:8888/metrics
- **Traces**: Distribuídos entre Service A → Service B

## 🧪 Testes

\`\`\`bash
# Testes unitários
make test

# Testes de integração
make test-integration

# Testes E2E
make test-e2e

# Cobertura de código
make coverage
\`\`\`

## 📚 Documentação

- [📖 Arquitetura](docs/ARCHITECTURE.md)
- [🔌 API Completa](docs/API.md)
- [⚙️ Desenvolvimento](docs/DEVELOPMENT.md)
- [🛠️ Monorepo Guide](docs/MONOREPO.md)
- [🔧 Troubleshooting](docs/TROUBLESHOOTING.md)
- [🤝 Contributing](docs/CONTRIBUTING.md)

## 🏃‍♂️ Comandos Úteis

\`\`\`bash
make help           # Ver todos os comandos
make up             # Subir ambiente
make down           # Parar ambiente
make logs           # Ver logs
make health         # Verificar saúde
make test-all       # Rodar todos os testes
\`\`\`
```

### OpenAPI Specification
```yaml
openapi: 3.0.0
info:
  title: Temperature by CEP API
  version: 1.0.0
  description: API para consulta de temperatura por CEP brasileiro

servers:
  - url: http://localhost:8080
    description: Service A - Input Service

paths:
  /:
    post:
      summary: Consulta temperatura por CEP
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - cep
              properties:
                cep:
                  type: string
                  pattern: '^[0-9]{8}$'
                  example: "01310100"
      responses:
        '200':
          description: Temperatura encontrada
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TemperatureResponse'
        '404':
          description: CEP não encontrado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '422':
          description: CEP inválido
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

components:
  schemas:
    TemperatureResponse:
      type: object
      properties:
        city:
          type: string
          example: "São Paulo"
        temp_C:
          type: number
          format: float
          example: 28.5
        temp_F:
          type: number
          format: float
          example: 83.3
        temp_K:
          type: number
          format: float
          example: 301.65
    
    ErrorResponse:
      type: object
      properties:
        message:
          type: string
          example: "invalid zipcode"
```

## Tarefas

### Documentação do Monorepo
- [ ] Criar README.md principal com estrutura do monorepo
- [ ] Adicionar docs/MONOREPO.md com guia específico
- [ ] Criar ARCHITECTURE.md com diagramas do sistema distribuído
- [ ] Escrever API.md com todos os endpoints
- [ ] Documentar processo de desenvolvimento no monorepo
- [ ] Criar guia de troubleshooting
- [ ] Adicionar CONTRIBUTING.md
- [ ] Documentar READMEs específicos de cada serviço

### Documentação da API
- [ ] Criar especificação OpenAPI completa
- [ ] Gerar documentação Swagger UI
- [ ] Criar coleção Postman
- [ ] Adicionar exemplos de requisições
- [ ] Documentar códigos de erro

### Testes Packages Compartilhados
- [ ] Testes para pkg/models
- [ ] Testes para pkg/telemetry
- [ ] Testes para pkg/utils
- [ ] Testes de integração entre packages
- [ ] Documentação dos packages compartilhados

### Testes Service A
- [ ] Testes unitários do validador
- [ ] Testes unitários do handler
- [ ] Testes do cliente HTTP
- [ ] Testes de integração da API
- [ ] Testes E2E com Service B mockado

### Testes Service B
- [ ] Testes do conversor de temperatura
- [ ] Testes dos clientes de API
- [ ] Testes do serviço de cache
- [ ] Mocks para APIs externas
- [ ] Testes de integração completos
- [ ] Benchmarks de performance

### Testes E2E do Sistema
- [ ] Setup do ambiente de testes E2E na raiz
- [ ] Testes de fluxo completo entre serviços
- [ ] Validação de traces no Zipkin
- [ ] Testes de carga do sistema completo
- [ ] Testes de resiliência
- [ ] Testes com testcontainers

### CI/CD Monorepo
- [ ] Configurar GitHub Actions para monorepo
- [ ] Pipeline de testes para cada serviço
- [ ] Pipeline de testes para packages compartilhados
- [ ] Análise de cobertura consolidada
- [ ] Build e push de imagens multi-service
- [ ] Deploy automático (opcional)

### Ferramentas e Scripts Monorepo
- [ ] Makefile completo na raiz
- [ ] Scripts de desenvolvimento em scripts/
- [ ] Script para rodar todos os testes
- [ ] Gerador de dados de teste
- [ ] Health check automatizado do sistema
- [ ] Script de análise de logs consolidado

## Critérios de Aceitação

1. Documentação completa e atualizada
2. Cobertura de testes > 80%
3. Todos os testes passando
4. API documentation acessível
5. Exemplos funcionais
6. CI/CD pipeline funcionando
7. Guias claros para novos desenvolvedores

## Exemplos de Testes

### Teste Unitário - Validador
```go
func TestCEPValidator(t *testing.T) {
    tests := []struct {
        name    string
        cep     string
        wantErr bool
    }{
        {"valid CEP", "12345678", false},
        {"invalid length", "1234567", true},
        {"with letters", "1234567a", true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateCEP(tt.cep)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateCEP() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Teste E2E
```go
func TestCompleteFlow(t *testing.T) {
    // Setup
    client := &http.Client{Timeout: 10 * time.Second}
    
    // Test valid CEP
    resp := postCEP(t, client, "01310100")
    assert.Equal(t, 200, resp.StatusCode)
    
    var result TemperatureResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    assert.NotEmpty(t, result.City)
    assert.Greater(t, result.TempC, -100.0)
    assert.Greater(t, result.TempF, -100.0)
    assert.Greater(t, result.TempK, 0.0)
    
    // Verify trace in Zipkin
    trace := getTraceFromZipkin(t, resp.Header.Get("X-Trace-Id"))
    assert.NotNil(t, trace)
    assert.Len(t, trace.Spans, 5) // Expected number of spans
}
```

## Makefile (Monorepo)
```makefile
.PHONY: help test docs build up down

# Monorepo settings
SERVICES := service-a service-b
PACKAGES := models telemetry utils

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Development commands
up: ## Start all services with docker-compose
	@echo "🚀 Starting all services..."
	@docker-compose up -d

down: ## Stop all services
	@echo "🛑 Stopping all services..."
	@docker-compose down

logs: ## Show logs from all services
	@docker-compose logs -f

# Testing commands
test: ## Run all unit tests
	@echo "🧪 Running unit tests for all services and packages..."
	@go test ./services/... ./pkg/... -v -cover

test-service-a: ## Run tests for service-a only
	@echo "🧪 Testing service-a..."
	@go test ./services/service-a/... -v -cover

test-service-b: ## Run tests for service-b only
	@echo "🧪 Testing service-b..."
	@go test ./services/service-b/... -v -cover

test-pkg: ## Run tests for shared packages
	@echo "🧪 Testing shared packages..."
	@go test ./pkg/... -v -cover

test-integration: ## Run integration tests
	@echo "🔗 Running integration tests..."
	@docker-compose up -d
	@sleep 10
	@go test ./services/*/test/integration/... -v
	@docker-compose down

test-e2e: ## Run E2E tests
	@echo "🌐 Running E2E tests..."
	@docker-compose up -d
	@sleep 15
	@go test ./test/e2e/... -v
	@docker-compose down

test-all: test test-integration test-e2e ## Run all tests

# Coverage
coverage: ## Generate coverage report for entire monorepo
	@echo "📊 Generating coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

coverage-service-a: ## Coverage for service-a
	@go test ./services/service-a/... -coverprofile=coverage-service-a.out
	@go tool cover -html=coverage-service-a.out -o coverage-service-a.html

coverage-service-b: ## Coverage for service-b
	@go test ./services/service-b/... -coverprofile=coverage-service-b.out
	@go tool cover -html=coverage-service-b.out -o coverage-service-b.html

# Build commands
build: ## Build all services
	@echo "🔨 Building all services..."
	@docker-compose build

build-service-a: ## Build service-a
	@echo "🔨 Building service-a..."
	@cd services/service-a && go build -o ../../bin/service-a ./cmd/server

build-service-b: ## Build service-b
	@echo "🔨 Building service-b..."
	@cd services/service-b && go build -o ../../bin/service-b ./cmd/server

# Documentation
docs: ## Generate API documentation
	@echo "📚 Generating API docs..."
	@swagger generate spec -o ./api/openapi.yaml --scan-models

# Code quality
lint: ## Run linters for entire monorepo
	@echo "🔍 Running linters..."
	@golangci-lint run ./...

lint-fix: ## Fix linting issues
	@golangci-lint run --fix ./...

fmt: ## Format code
	@go fmt ./...

# Health checks
health: ## Check health of all services
	@echo "🩺 Checking service health..."
	@curl -f http://localhost:8080/health || echo "Service A unhealthy"
	@curl -f http://localhost:8081/health || echo "Service B unhealthy"

# Monitoring
zipkin: ## Open Zipkin UI
	@open http://localhost:9411

# Utilities
clean: ## Clean build artifacts and containers
	@echo "🧹 Cleaning up..."
	@docker-compose down -v
	@docker system prune -f
	@rm -rf bin/ coverage*.out coverage*.html

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
```

## GitHub Actions Workflow (Monorepo)
```yaml
name: CI - Monorepo

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: '1.21'

jobs:
  # Detectar mudanças para otimizar builds
  changes:
    runs-on: ubuntu-latest
    outputs:
      service-a: ${{ steps.changes.outputs.service-a }}
      service-b: ${{ steps.changes.outputs.service-b }}
      pkg: ${{ steps.changes.outputs.pkg }}
      docs: ${{ steps.changes.outputs.docs }}
    steps:
    - uses: actions/checkout@v3
    - uses: dorny/paths-filter@v2
      id: changes
      with:
        filters: |
          service-a:
            - 'services/service-a/**'
            - 'pkg/**'
            - 'go.mod'
            - 'go.sum'
          service-b:
            - 'services/service-b/**'
            - 'pkg/**'
            - 'go.mod'
            - 'go.sum'
          pkg:
            - 'pkg/**'
            - 'go.mod'
            - 'go.sum'
          docs:
            - 'docs/**'
            - '*.md'

  # Testes dos packages compartilhados
  test-pkg:
    runs-on: ubuntu-latest
    needs: changes
    if: needs.changes.outputs.pkg == 'true'
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Test shared packages
      run: make test-pkg
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: pkg

  # Testes do Service A
  test-service-a:
    runs-on: ubuntu-latest
    needs: changes
    if: needs.changes.outputs.service-a == 'true'
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Test Service A
      run: make test-service-a
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage-service-a.out
        flags: service-a

  # Testes do Service B
  test-service-b:
    runs-on: ubuntu-latest
    needs: changes
    if: needs.changes.outputs.service-b == 'true'
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Test Service B
      run: make test-service-b
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage-service-b.out
        flags: service-b

  # Build e testes E2E
  test-e2e:
    runs-on: ubuntu-latest
    needs: [test-pkg, test-service-a, test-service-b]
    if: always() && (needs.test-pkg.result == 'success' || needs.test-pkg.result == 'skipped') && (needs.test-service-a.result == 'success' || needs.test-service-a.result == 'skipped') && (needs.test-service-b.result == 'success' || needs.test-service-b.result == 'skipped')
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Build images
      run: make build
    
    - name: Run E2E tests
      run: make test-e2e
    
    - name: Upload E2E coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: e2e

  # Linting
  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
        args: --timeout=5m
```

## Referências
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Testcontainers Go](https://golang.testcontainers.org/)
- [GitHub Actions for Go](https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go)