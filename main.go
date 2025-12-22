package main

import (
	"go-microservice/internal/handlers"
	"go-microservice/internal/metrics"
	"go-microservice/internal/services"
	"go-microservice/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	// Prometheus metrics
	metrics.Init()

	// Minio client
	minioClient, err := utils.NewMinioClient(
		"minio:9000",
		"minioadmin",
		"minioadmin",
		"audit-logs",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Integration service
	integrationService := services.NewIntegrationService(minioClient)

	// Async components
	auditLogger := utils.NewAuditLogger(1000, integrationService)
	notifier := utils.NewNotifier(1000)

	auditLogger.Start()
	notifier.Start()

	// Services
	userService := services.NewUserService()

	// Handlers
	userHandler := handlers.NewUserHandler(
		userService,
		auditLogger,
		notifier,
	)

	integrationHandler := handlers.NewIntegrationHandler()

	// Router
	r := mux.NewRouter()

	// API router
	api := r.PathPrefix("/api").Subrouter()

	// Rate limiter
	rateLimiter := utils.NewRateLimiter(1000, 1000)
	api.Use(rateLimiter.Middleware)

	// Metrics middleware
	api.Use(metrics.Middleware)

	// User routes
	api.HandleFunc("/users", userHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", userHandler.Get).Methods(http.MethodGet)
	api.HandleFunc("/users", userHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", userHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", userHandler.Delete).Methods(http.MethodDelete)

	// Prometheus endpoint
	r.HandleFunc("/metrics", integrationHandler.Metrics).Methods(http.MethodGet)

	// HTTP server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("HTTP server started on :8080")
	log.Fatal(server.ListenAndServe())
}
