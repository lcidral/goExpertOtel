# PRP-001: Implementação do Serviço A - Input Service

## Resumo
Implementar o serviço responsável por receber e validar CEPs, encaminhando requisições válidas para o Serviço B.

## Motivação
O Serviço A atua como gateway de entrada do sistema, responsável pela validação inicial dos dados e roteamento para o serviço de orquestração. Sua implementação é fundamental para garantir que apenas dados válidos sejam processados pelo sistema.

## Descrição Detalhada

### Arquitetura
- Servidor HTTP em Go
- Porta padrão: 8080
- Comunicação síncrona com Serviço B via HTTP

### Funcionalidades
1. **Endpoint POST /**
   - Recebe JSON: `{ "cep": "29902555" }`
   - Valida formato do CEP (8 dígitos, tipo string)
   - Encaminha para Serviço B se válido
   - Retorna erro 422 se inválido

2. **Validações**
   - CEP deve ser string
   - CEP deve conter exatamente 8 dígitos
   - Apenas números são permitidos

3. **Integração com Serviço B**
   - Cliente HTTP configurável via variável de ambiente
   - Timeout configurável
   - Retry policy para resiliência

## Implementação Proposta

### Estrutura de Monorepo
Este projeto utiliza uma estrutura de monorepo para facilitar o desenvolvimento, compartilhamento de código e gerenciamento de dependências entre os serviços.

```
goExpertOtel/                    # Raiz do monorepo
├── services/
│   ├── service-a/              # Serviço A - Input Service
│   │   ├── cmd/
│   │   │   └── server/
│   │   │       └── main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   │   └── cep_handler.go
│   │   │   ├── validator/
│   │   │   │   └── cep_validator.go
│   │   │   └── client/
│   │   │       └── service_b_client.go
│   │   └── config/
│   │       └── config.go
│   └── service-b/              # Serviço B (implementado no PRP-002)
├── pkg/                        # Código compartilhado
│   ├── models/
│   │   └── cep.go
│   └── telemetry/              # Compartilhado no PRP-003
├── go.mod                      # Módulo Go único para o monorepo
├── go.sum
├── go.work                     # Go workspace (opcional)
├── docker-compose.yml
├── Makefile
└── README.md
```

### Module Path
O módulo Go principal será: `github.com/yourusername/goExpertOtel`

Imports internos do Service A:
```go
import (
    "github.com/yourusername/goExpertOtel/pkg/models"
    "github.com/yourusername/goExpertOtel/pkg/telemetry"
)
```

### Dependências
- `net/http` - Servidor HTTP
- `encoding/json` - Parsing JSON
- `github.com/go-chi/chi/v5` - Router HTTP
- `github.com/joho/godotenv` - Variáveis de ambiente

### Configurações
```env
SERVICE_B_URL=http://service-b:8081
REQUEST_TIMEOUT=30s
PORT=8080
```

## Tarefas

- [ ] Inicializar monorepo com go.mod na raiz
- [ ] Criar estrutura de diretórios services/service-a
- [ ] Implementar modelo compartilhado em pkg/models/cep.go
- [ ] Criar validador de CEP com regex
- [ ] Implementar handler HTTP para receber requisições
- [ ] Criar cliente HTTP para comunicação com Serviço B
- [ ] Implementar middleware de logging
- [ ] Adicionar health check endpoint
- [ ] Criar testes unitários para validador
- [ ] Criar testes de integração para handler
- [ ] Implementar graceful shutdown
- [ ] Adicionar métricas básicas (requisições/segundo, erros)
- [ ] Configurar Dockerfile com contexto do monorepo

## Critérios de Aceitação

1. Serviço responde corretamente a requisições válidas
2. Retorna 422 com mensagem "invalid zipcode" para CEPs inválidos
3. Propaga corretamente a resposta do Serviço B
4. Testes com cobertura mínima de 80%
5. Logs estruturados para debugging
6. Documentação da API

## Riscos e Mitigações

| Risco | Mitigação |
|-------|-----------|
| Timeout na comunicação com Serviço B | Implementar circuit breaker |
| CEPs com formato especial | Normalizar entrada removendo caracteres não numéricos |
| Alta carga de requisições | Implementar rate limiting |

## Referências
- [Go HTTP Server Best Practices](https://golang.org/doc/articles/wiki/)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [Go Project Layout](https://github.com/golang-standards/project-layout)