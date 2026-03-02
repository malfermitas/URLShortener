package webui

import (
	"html/template"
	"net/http"
	"urlshortener/internal/logging"
	"urlshortener/internal/metrics"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	templates *template.Template
}

func NewHandler(templatesDir string) *Handler {
	templates := template.Must(template.ParseGlob(templatesDir + "/*.html"))
	return &Handler{
		templates: templates,
	}
}

func (h *Handler) ServeHTML(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Status(http.StatusOK)

	if err := h.templates.ExecuteTemplate(ctx.Writer, "index.html", nil); err != nil {
		logging.AppLogger.Error("Failed to render template", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	logging.AppLogger.Debug("Served index.html")
	if metrics.WebUIPageviewsTotal != nil {
		metrics.WebUIPageviewsTotal.Inc()
	}
}

// ServeAnalyticsHTML serves the analytics UI page for viewing analytics.
func (h *Handler) ServeAnalyticsHTML(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Status(http.StatusOK)

	if err := h.templates.ExecuteTemplate(ctx.Writer, "analytics.html", nil); err != nil {
		logging.AppLogger.Error("Failed to render analytics template", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	logging.AppLogger.Debug("Served analytics.html")
	if metrics.WebUIPageviewsTotal != nil {
		metrics.WebUIPageviewsTotal.Inc()
	}
}
