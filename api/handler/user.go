package handler

import (
	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func AddUser(c echo.Context, queries *out.Queries) (pgtype.UUID, error) {
	id, err := queries.CreateUserAndReturnId(c.Request().Context(), "HelloWorld")

	return id, err
}
