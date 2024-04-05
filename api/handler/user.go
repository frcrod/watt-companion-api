package handler

import (
	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/frcrod/watt-companion-api/web"
	"github.com/frcrod/watt-companion-api/web/pages"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func AddUser(c echo.Context, queries *out.Queries) (pgtype.UUID, error) {
	id, err := queries.CreateUserAndReturnId(c.Request().Context(), out.CreateUserAndReturnIdParams{
		Email:    "Email",
		Nickname: "email",
	})

	return id, err
}

func Home(c echo.Context) error {
	return web.Render(c, pages.Home())
}
