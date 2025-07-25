package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	apiVersionV1 = "/v1"
)

type Router struct {
	echo       *echo.Echo
	nftHandler *NFTHandler
}

func NewRouter(echo *echo.Echo, nftHandler *NFTHandler) *Router {
	return &Router{
		echo:       echo,
		nftHandler: nftHandler,
	}
}

func (r *Router) RegisterRoutes() {
	r.echo.Use(middleware.Logger())

	apiVersionV1 := r.echo.Group(apiVersionV1)

	r.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	apiVersionV1.GET("/ownerships/:id", r.nftHandler.GetOwnership)
	apiVersionV1.GET("/trait-rarities", r.nftHandler.GetTraitRarities)

}
