# Service A - Input Service

Serviço responsável por receber e validar CEPs, encaminhando requisições válidas para o Serviço B.

## Funcionalidades

- **Validação de CEP**: Valida formato brasileiro (8 dígitos)
- **Normalização**: Remove caracteres especiais automaticamente
- **Proxy para Serviço B**: Encaminha requisições válidas
- **Health Check**: Endpoint de monitoramento
- **Graceful Shutdown**: Encerramento seguro do serviço

## Endpoints

### POST /
Recebe e processa CEPs para consulta de temperatura.

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
  "temp_K": 298.65
}
```

**Error Responses:**
- `422`: CEP inválido (`{"message": "invalid zipcode"}`)
- `404`: CEP não encontrado (`{"message": "can not find zipcode"}`)

### GET /health
Endpoint de health check.

**Response (200):**
```json
{
  "status": "healthy",
  "service": "service-a",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Configuração

### Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PORT` | `8080` | Porta do serviço |
| `SERVICE_B_URL` | `http://localhost:8081` | URL do Serviço B |
| `REQUEST_TIMEOUT` | `30s` | Timeout para chamadas ao Serviço B |

## Execução

### Desenvolvimento Local
```bash
# Da raiz do monorepo
cd services/service-a
go run cmd/server/main.go
```

### Build
```bash
# Da raiz do monorepo
go build -o bin/service-a ./services/service-a/cmd/server
./bin/service-a
```

## Testes

```bash
# Testes unitários
go test ./services/service-a/... -v

# Testes com cobertura
go test ./services/service-a/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Exemplos de Uso

```bash
# CEP válido
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# CEP com formatação (será normalizado)
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310-100"}'

# Health check
curl http://localhost:8080/health
```

## Arquitetura

```
services/service-a/
├── cmd/server/           # Aplicação principal
├── internal/
│   ├── handler/         # Handlers HTTP
│   ├── validator/       # Validação de CEP
│   └── client/          # Cliente do Serviço B
└── config/              # Configurações
```

## Dependências

- `github.com/go-chi/chi/v5` - Router HTTP
- `github.com/joho/godotenv` - Variáveis de ambiente
- Compartilha `pkg/models` com outros serviços