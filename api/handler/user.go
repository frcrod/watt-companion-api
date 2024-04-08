package handler

import (
	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/frcrod/watt-companion-api/internal/types"
	"github.com/frcrod/watt-companion-api/web"
	"github.com/frcrod/watt-companion-api/web/component"
	"github.com/frcrod/watt-companion-api/web/page"
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

func newAppliance(name string, wattage float32) types.Appliance {
	return types.Appliance{
		Name:    name,
		Wattage: wattage,
	}
}

func Home(c echo.Context) error {
	appliances := []types.Appliance{
		newAppliance("Hello", 32.32),
		newAppliance("Hello", 32.32),
	}

	props := component.ApplianceProp{
		Appliances: appliances,
	}

	return web.Render(c, page.Home(props))
}
