package handler

import (
	"errors"
	"net/http"
	"urlshortener/internal/core/port/in"

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
	originalURL, err := r.urlService.GetOriginal(ctx, shortKey)
	if err != nil {
		if errors.Is(err, in.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, err.Error())
		} else {
			ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Redirect(http.StatusMovedPermanently, originalURL)
}
