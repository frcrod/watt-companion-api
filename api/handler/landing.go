package handler

import (
	"github.com/frcrod/watt-companion-api/web"
	"github.com/frcrod/watt-companion-api/web/component"
	"github.com/labstack/echo/v4"
)

func HandleLanding(c echo.Context) error {
	return web.Render(c, component.Nav())
}
