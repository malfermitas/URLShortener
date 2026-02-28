package handler

import (
	"errors"
	"net/http"
	"time"
	"urlshortener/internal/core/model"
	"urlshortener/internal/core/port/in"
	"urlshortener/internal/logging"

	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/ginext"
)

type RedirectHandler interface {
	Redirect(ctx *ginext.Context)
}

type redirectHandler struct {
	urlService in.URLService
}

func NewRedirectHandler(urlService in.URLService) RedirectHandler {
	return &redirectHandler{
		urlService: urlService,
	}
}

func (r redirectHandler) Redirect(ctx *ginext.Context) {
	shortKey := ctx.Param("short_url")

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Var(shortKey, "required,alphanum"); err != nil {
		logging.AppLogger.Debug("Invalid short URL format", "short_key", shortKey, "error", err.Error())
		ctx.JSON(http.StatusBadRequest, "invalid short url format")
		return
	}

	originalURL, err := r.urlService.GetOriginal(ctx, shortKey)
	if err != nil {
		if errors.Is(err, in.ErrNotFound) {
			logging.AppLogger.Debug("Short URL not found", "short_key", shortKey)
			ctx.JSON(http.StatusNotFound, err.Error())
		} else {
			logging.AppLogger.Error("Failed to get original URL", err, "short_key", shortKey)
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	go func(shortKey string, userAgent string, ip string, referrer string) {
		err := r.urlService.RecordHit(ctx.Request.Context(), &model.URLHitEvent{
			URLID:     shortKey,
			UserAgent: userAgent,
			IP:        ip,
			Referrer:  referrer,
			Timestamp: time.Now(),
		})
		if err != nil {
			logging.AppLogger.Error("failed to record hit: %v", err)
		}
	}(shortKey, ctx.GetHeader("User-Agent"), ctx.ClientIP(), ctx.GetHeader("Referer"))

	logging.AppLogger.Debug("Redirecting", "short_key", shortKey, "original_url", originalURL)
	ctx.Redirect(http.StatusMovedPermanently, originalURL)
}
