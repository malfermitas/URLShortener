package rest

import (
	"fmt"
	"urlshortener/internal/adapter/in/rest/handler"
	restmiddleware "urlshortener/internal/adapter/in/rest/middleware"
	"urlshortener/internal/adapter/in/webui"
	"urlshortener/internal/metrics"
	"urlshortener/internal/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(r handler.RedirectHandler, s handler.ShortenerHandler, a handler.AnalyticsHandler, w *webui.Handler, templatesDir string) *gin.Engine {
	router := gin.Default()

	router.Use(otelgin.Middleware("urlshortener"))
	router.Use(panicRecoveryMiddleware())
	router.Use(restmiddleware.GinLogger())
	router.Use(metrics.MetricsMiddleware())
	router.LoadHTMLGlob(templatesDir + "/*.html")

	router.GET("/s/:short_url", r.Redirect)
	router.GET("/s/:short_url/analytics", a.GetAnalytics)
	router.POST("/shorten", s.Shorten)
	router.GET("/", w.ServeHTML)
	// WebUI analytics page
	router.GET("/analytics", w.ServeAnalyticsHTML)

	return router
}

func panicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx, span := tracing.StartSpan(c.Request.Context(), "panic")
				tracing.RecordError(ctx, fmt.Errorf("%v", err))
				span.End()

				c.AbortWithStatusJSON(500, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
