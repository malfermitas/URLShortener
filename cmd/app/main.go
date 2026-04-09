package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"urlshortener/internal/adapter/in/rest"
	"urlshortener/internal/adapter/in/rest/handler"
	"urlshortener/internal/adapter/in/webui"
	"urlshortener/internal/adapter/out/generator"
	"urlshortener/internal/adapter/out/postgres"
	"urlshortener/internal/adapter/out/redis"
	"urlshortener/internal/config"
	"urlshortener/internal/core/service"
	"urlshortener/internal/logging"
	"urlshortener/internal/metrics"
	"urlshortener/internal/tracing"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logging.AppLogger, err = logging.NewURLShortenerLogger()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	if err := tracing.InitTracing(cfg.Tracing); err != nil {
		logging.AppLogger.Warn("Failed to initialize tracing", "error", err.Error())
	}

	templatesDir := "internal/adapter/in/webui"

	// Initialize Prometheus metrics
	metrics.InitMetrics()

	urlRepo, err := postgres.NewURLRepository(cfg.Database.DSN())
	if err != nil {
		logging.AppLogger.Error("Cannot connect to postgres database", err)
		os.Exit(1)
	}

	hitEventRepo, err := postgres.NewURLHitEventRepository(cfg.Database.DSN())
	if err != nil {
		logging.AppLogger.Error("Cannot connect to postgres database", err)
		os.Exit(1)
	}

	urlCache, err := redis.NewURLCache(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.MaxRecordTTL)
	if err != nil {
		logging.AppLogger.Error("Cannot connect to Redis", err)
		os.Exit(1)
	}
	logging.AppLogger.Info("Connected to Redis")

	keyGenerator := generator.NewURLGenerator()

	urlService := service.NewUrlService(urlRepo, keyGenerator, hitEventRepo, urlCache)

	shortenerHandler := handler.NewShortenerHandler(urlService)
	redirectHandler := handler.NewRedirectHandler(urlService)
	analyticsHandler := handler.NewAnalyticsHandler(urlService)
	webUIHandler := webui.NewHandler(templatesDir)

	router := rest.NewRouter(redirectHandler, shortenerHandler, analyticsHandler, webUIHandler, templatesDir)
	// Apply Prometheus metrics middleware for HTTP metrics
	router.Use(metrics.MetricsMiddleware())
	// Expose Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Pprof endpoints
	router.GET("/debug/pprof/*path", gin.WrapH(http.DefaultServeMux))
	router.GET("/debug/pprof", gin.WrapH(http.DefaultServeMux))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logging.AppLogger.Info(fmt.Sprintf("Starting server on %s", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.AppLogger.Error("Server error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.AppLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.AppLogger.Error("Server forced to shutdown", err)
	}

	logging.AppLogger.Info("Server exited")
}
