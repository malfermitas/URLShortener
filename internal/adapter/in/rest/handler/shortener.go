package handler

import (
	"net/http"
	"urlshortener/internal/core/port/in"
	"urlshortener/internal/logging"
	"urlshortener/internal/tracing"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type ShortenerHandler interface {
	Shorten(ctx *ginext.Context)
}

type shortenerHandler struct {
	service in.URLService
}

type ShortenRequest struct {
	OriginURL string `json:"origin_url" binding:"required,url"`
	CustomURL string `json:"custom_url,omitzero" binding:"omitempty,alphanum,min=3,max=20"`
}

func NewShortenerHandler(urlService in.URLService) ShortenerHandler {
	return &shortenerHandler{
		service: urlService,
	}
}

func (s shortenerHandler) Shorten(ctx *ginext.Context) {
	var req ShortenRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logging.AppLogger.Error(
			"shorten request bind error: %s", err.Error(),
			"trace_id", tracing.GetTraceID(ctx.Request.Context()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortURL, err := s.service.Create(ctx.Request.Context(), req.OriginURL, req.CustomURL)
	if err != nil {
		logging.AppLogger.Error(
			"shorten service create error: %s", err.Error(),
			"trace_id", tracing.GetTraceID(ctx.Request.Context()),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logging.AppLogger.Info(
		"URL shortened successfully",
		"short_url", shortURL,
		"trace_id", tracing.GetTraceID(ctx.Request.Context()),
	)
	ctx.JSON(http.StatusCreated, gin.H{"short_url": shortURL})
}
