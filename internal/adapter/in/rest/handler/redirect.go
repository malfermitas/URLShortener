package handler

import (
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
	shortKey := ctx.Param("id")
	originalURL, err := r.urlService.GetOriginal(ctx, shortKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	ctx.Redirect(http.StatusMovedPermanently, originalURL)
}
