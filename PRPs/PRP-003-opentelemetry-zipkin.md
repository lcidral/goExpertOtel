# PRP-003: Implementação de OpenTelemetry e Zipkin

## Resumo
Adicionar observabilidade completa ao sistema através da implementação de tracing distribuído com OpenTelemetry e visualização com Zipkin, permitindo monitoramento e debugging eficaz das requisições entre os serviços.

## Motivação
Em sistemas distribuídos, a observabilidade é fundamental para entender o comportamento das aplicações, identificar gargalos de performance e diagnosticar problemas. OpenTelemetry fornece um padrão unificado para coleta de telemetria, enquanto Zipkin permite visualização intuitiva dos traces.

## Descrição Detalhada

### Componentes de Observabilidade

1. **OpenTelemetry SDK**
   - Instrumentação automática para HTTP
   - Criação manual de spans customizados
   - Propagação de contexto entre serviços
   - Exportação de traces para collector

2. **OTEL Collector**
   - Recepção de traces dos serviços
   - Processamento e enriquecimento
   - Exportação para Zipkin
   - Configuração de pipelines

3. **Zipkin**
   - Interface web para visualização
   - Análise de latência
   - Dependency graph
   - Busca e filtros avançados

### Spans a Implementar

#### Service A
- `http.request` - Span principal da requisição
- `cep.validation` - Validação do CEP
- `service_b.call` - Chamada para Service B

#### Service B
- `http.request` - Span principal da requisição  
- `cep.validation` - Revalidação do CEP
- `viacep.api.call` - Chamada para ViaCEP
- `weather.api.call` - Chamada para WeatherAPI
- `temperature.conversion` - Conversão de temperaturas
- `cache.lookup` - Consulta ao cache
- `cache.store` - Armazenamento no cache

## Implementação Proposta

### Estrutura no Monorepo
A implementação de OpenTelemetry será compartilhada entre todos os serviços através do diretório `pkg/telemetry`.

```
goExpertOtel/                    # Raiz do monorepo
├── services/
│   ├── service-a/              # Utilizará pkg/telemetry
│   └── service-b/              # Utilizará pkg/telemetry
├── pkg/
│   ├── telemetry/              # Código de telemetria compartilhado
│   │   ├── tracer.go          # Inicialização do tracer
│   │   ├── middleware.go      # Middleware HTTP com tracing
│   │   ├── propagation.go     # Propagação de contexto
│   │   └── config.go          # Configurações OTEL
│   └── models/
├── deployments/
│   ├── otel-collector/
│   │   ├── otel-collector-config.yaml
│   │   └── Dockerfile
│   └── docker/
│       └── docker-compose.yml
├── go.mod                      # Módulo único do monorepo
└── go.sum
```

### Dependências OpenTelemetry
```go
// go.mod additions
require (
    go.opentelemetry.io/otel v1.19.0
    go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.19.0
    go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.19.0
    go.opentelemetry.io/otel/sdk v1.19.0
    go.opentelemetry.io/otel/trace v1.19.0
    go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.45.0
)
```

### Configuração OTEL Collector
```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
  attributes:
    actions:
      - key: environment
        value: development
        action: insert

exporters:
  zipkin:
    endpoint: "http://zipkin:9411/api/v2/spans"
  logging:
    loglevel: debug

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch, attributes]
      exporters: [zipkin, logging]
```

### Código de Inicialização do Tracer
```go
// pkg/telemetry/tracer.go
package telemetry

import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
    ctx := context.Background()
    
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint("otel-collector:4317"),
        otlptracegrpc.WithInsecure(),
    )
    
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String("1.0.0"),
        )),
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
    )
    
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.TraceContext{})
    
    return tp, nil
}
```

### Uso nos Serviços
```go
// services/service-a/cmd/server/main.go
import "github.com/yourusername/goExpertOtel/pkg/telemetry"

func main() {
    tp, err := telemetry.InitTracer("service-a")
    if err != nil {
        log.Fatal(err)
    }
    defer tp.Shutdown(context.Background())
    
    // resto da inicialização...
}
```

## Tarefas

### Setup Inicial no Monorepo
- [ ] Criar diretório pkg/telemetry no monorepo
- [ ] Implementar tracer.go com inicialização compartilhada
- [ ] Criar middleware.go para auto-instrumentação HTTP
- [ ] Implementar propagation.go para contexto distribuído
- [ ] Adicionar dependências OTEL ao go.mod principal

### Service A
- [ ] Importar pkg/telemetry no service-a
- [ ] Instrumentar main.go com tracer compartilhado
- [ ] Adicionar middleware de tracing nos handlers
- [ ] Criar span para validação de CEP
- [ ] Instrumentar cliente HTTP para Service B
- [ ] Adicionar atributos customizados aos spans

### Service B  
- [ ] Importar pkg/telemetry no service-b
- [ ] Instrumentar main.go com tracer compartilhado
- [ ] Adicionar middleware de tracing
- [ ] Criar spans para chamadas às APIs externas
- [ ] Instrumentar operações de cache
- [ ] Adicionar span para conversões
- [ ] Enriquecer spans com dados de erro

### Infraestrutura
- [ ] Criar diretório deployments/otel-collector
- [ ] Configurar OTEL Collector em deployments/
- [ ] Adicionar Zipkin ao docker-compose principal
- [ ] Configurar pipelines de processamento
- [ ] Criar dashboards customizados no Zipkin
- [ ] Documentar queries úteis

### Monitoramento
- [ ] Implementar health checks para collector
- [ ] Adicionar métricas de performance
- [ ] Configurar alertas básicos
- [ ] Criar runbook de troubleshooting

## Critérios de Aceitação

1. Traces completos de ponta a ponta visíveis no Zipkin
2. Propagação correta do trace context entre serviços
3. Spans com nomes e atributos significativos
4. Latência de cada operação claramente identificável
5. Erros e exceções capturados nos traces
6. Performance overhead < 5%
7. Documentação de uso e análise

## Configurações Adicionais

### Variáveis de Ambiente
```env
# Service A & B
OTEL_SERVICE_NAME=service-a
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
OTEL_EXPORTER_OTLP_INSECURE=true
OTEL_TRACES_SAMPLER=always_on
OTEL_METRICS_EXPORTER=none
OTEL_LOGS_EXPORTER=none

# Collector
OTEL_COLLECTOR_ZIPKIN_ENDPOINT=http://zipkin:9411/api/v2/spans
```

## Riscos e Mitigações

| Risco | Mitigação |
|-------|-----------|
| Performance overhead | Implementar sampling adaptativo |
| Volume alto de dados | Configurar retention policies |
| Perda de traces | Implementar buffer persistente |
| Collector indisponível | Fallback para exportação direta |

## Melhores Práticas

1. **Naming Convention**
   - Spans: `<component>.<operation>`
   - Atributos: snake_case
   - Recursos: seguir semantic conventions

2. **Atributos Essenciais**
   - `http.method`, `http.url`, `http.status_code`
   - `cep.value`, `city.name`
   - `cache.hit`, `api.name`
   - `error`, `error.message`

3. **Performance**
   - Usar batch export
   - Limitar atributos por span
   - Implementar sampling inteligente

## Exemplos de Queries Zipkin

```
# Buscar traces com erros
error=true

# Traces lentos (> 1s)
minDuration=1000ms

# Por CEP específico
tagQuery=cep.value="12345678"

# Falhas na WeatherAPI
serviceName=service-b AND tagQuery=api.name="weather" AND error=true
```

## Referências
- [OpenTelemetry Go Getting Started](https://opentelemetry.io/docs/languages/go/getting-started/)
- [OpenTelemetry Collector Configuration](https://opentelemetry.io/docs/collector/configuration/)
- [Zipkin Documentation](https://zipkin.io/pages/documentation.html)
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)
- [OpenTelemetry Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/)