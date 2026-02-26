package handler

import (
	"net/http"
	"urlshortener/internal/core/port/in"
	"urlshortener/internal/logging"

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
	OriginURL string `json:"origin_url"`
	CustomURL string `json:"custom_url,omitempty"`
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
		logging.AppLogger.Error("shorten request bind error: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	shortURL, err := s.service.Create(ctx, req.OriginURL, req.CustomURL)
	if err != nil {
		logging.AppLogger.Error("shorten service create error: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"short_url": shortURL})
}
