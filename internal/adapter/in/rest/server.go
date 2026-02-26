package rest

import (
	"urlshortener/internal/adapter/in/rest/handler"

	"github.com/wb-go/wbf/ginext"
)

func NewRouter(r handler.RedirectHandler, s handler.ShortenerHandler) *ginext.Engine {
	router := ginext.New("")
	router.Use(ginext.Recovery())

	router.GET("/s/:short_url", r.Redirect)
	router.POST("/shorten", s.Shorten)

	return router
}
