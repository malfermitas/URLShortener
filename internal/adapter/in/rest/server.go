package rest

import (
	"urlshortener/internal/adapter/in/rest/handler"
	restmiddleware "urlshortener/internal/adapter/in/rest/middleware"
	"urlshortener/internal/adapter/in/webui"
	"urlshortener/internal/metrics"

	"github.com/gin-gonic/gin"
)

func NewRouter(r handler.RedirectHandler, s handler.ShortenerHandler, a handler.AnalyticsHandler, w *webui.Handler, templatesDir string) *gin.Engine {
	router := gin.Default()

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
