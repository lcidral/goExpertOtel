package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/lcidral/goExpertOtel/pkg/telemetry"
	"github.com/lcidral/goExpertOtel/services/service-b/config"
	"github.com/lcidral/goExpertOtel/services/service-b/internal/cache"
	"github.com/lcidral/goExpertOtel/services/service-b/internal/client"
	"github.com/lcidral/goExpertOtel/services/service-b/internal/handler"
)

func main() {
	// Load environment variables from .env file (if exists)
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: arquivo .env não encontrado: %v", err)
	}

	// Initialize OpenTelemetry
	tp, err := telemetry.InitTracer("service-b")
	if err != nil {
		log.Printf("Aviso: falha ao inicializar OpenTelemetry: %v", err)
	} else {
		defer func() {
			if err := telemetry.Shutdown(tp); err != nil {
				log.Printf("Erro ao finalizar OpenTelemetry: %v", err)
			}
		}()
		log.Println("🔍 OpenTelemetry inicializado com sucesso")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Validate configuration
	if err := cfg.ValidateConfig(); err != nil {
		log.Fatalf("Erro de configuração: %v", err)
	}

	// Initialize cache
	memoryCache := cache.NewMemoryCache(cfg.CacheTTL, cfg.CacheCleanup)

	// Initialize clients
	openCEPClient := client.NewOpenCEPClient(cfg.OpenCEPURL, cfg.RequestTimeout)
	weatherClient := client.NewWeatherClient(cfg.WeatherAPIURL, cfg.WeatherAPIKey, cfg.RequestTimeout)

	// Initialize handlers
	tempHandler := handler.NewTemperatureHandler(openCEPClient, weatherClient, memoryCache)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	
	// OpenTelemetry middleware
	if tp != nil {
		r.Use(telemetry.HTTPMiddleware("service-b"))
		log.Println("🔍 Middleware OpenTelemetry ativado")
	}

	// CORS for development
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	// Routes
	r.Post("/temperature", tempHandler.HandleTemperature)
	r.Get("/health", tempHandler.HealthCheck)
	r.Get("/cache/stats", tempHandler.CacheStats)

	// Server configuration
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌡️ Service B iniciado na porta %s", cfg.Port)
		log.Printf("📍 Health check: http://localhost:%s/health", cfg.Port)
		log.Printf("🌐 OpenCEP URL: %s", cfg.OpenCEPURL)
		log.Printf("☁️ WeatherAPI URL: %s", cfg.WeatherAPIURL)
		log.Printf("🗄️ Cache TTL: %v", cfg.CacheTTL)

		if cfg.WeatherAPIKey == "" {
			log.Printf("⚠️ ATENÇÃO: WEATHER_API_KEY não configurada!")
		}
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until signal is received
	<-c
	log.Println("🛑 Iniciando graceful shutdown...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Erro durante shutdown: %v", err)
	}

	// Clear cache
	memoryCache.Clear()
	log.Println("🗑️ Cache limpo")

	log.Println("✅ Service B encerrado com sucesso")
}