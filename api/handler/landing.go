package handler

import (
	"fmt"
	"log/slog"

	"github.com/frcrod/watt-companion-api/web"
	"github.com/frcrod/watt-companion-api/web/component"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func HandleLanding(e echo.Context) error {
	session, err := gothic.Store.Get(e.Request(), "auth")
	if err != nil {
		slog.Error("Error getting user session")
		return err
	}

	user := session.Values["user"]

	fmt.Println(user, session.Values["user_id"])

	return web.Render(e, component.Nav())
}
