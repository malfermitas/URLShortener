package handler

import (
	"errors"
	"net/http"
	"urlshortener/internal/core/port/in"
	"urlshortener/internal/logging"

	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/ginext"
)

type AnalyticsHandler interface {
	GetAnalytics(ctx *ginext.Context)
}

type analyticsHandler struct {
	urlService in.URLService
}

func NewAnalyticsHandler(urlService in.URLService) AnalyticsHandler {
	return &analyticsHandler{
		urlService: urlService,
	}
}

func (h analyticsHandler) GetAnalytics(ctx *ginext.Context) {
	shortKey := ctx.Param("short_url")

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Var(shortKey, "required,alphanum"); err != nil {
		logging.AppLogger.Debug("Invalid short URL format in analytics request", "short_key", shortKey, "error", err.Error())
		ctx.JSON(http.StatusBadRequest, "invalid short url format")
		return
	}

	analytics, err := h.urlService.GetAnalytics(ctx, shortKey)
	if err != nil {
		if errors.Is(err, in.ErrNotFound) || err.Error() == "URL not found" {
			logging.AppLogger.Debug("Analytics not found", "short_key", shortKey)
			ctx.JSON(http.StatusNotFound, "URL not found")
		} else {
			logging.AppLogger.Error("Failed to get analytics", err, "short_key", shortKey)
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	logging.AppLogger.Debug("Analytics request successful", "short_key", shortKey)
	ctx.JSON(http.StatusOK, analytics)
}
