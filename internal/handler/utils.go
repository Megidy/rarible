package handler

import "github.com/labstack/echo/v4"

func getFromParam(ctx echo.Context, param string) string {
	return ctx.Param(param)
}
