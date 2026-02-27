package webui

import (
	"html/template"
	"net/http"

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
		ctx.Status(http.StatusInternalServerError)
	}
}
