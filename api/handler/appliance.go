package handler

import (
	"math/big"

	"github.com/frcrod/watt-companion-api/api/util"
	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func AddAppliance(c echo.Context, queries *out.Queries) (pgtype.UUID, error) {
	uu, err := util.ConvertStringTo16Bytes("c3018885a487467cb1f794dc72850ab1")
	if err != nil {
		c.Logger().Fatal("ERROR CONVERTING STRING TO 16 BYTE SLICE")
	}

	wattage := pgtype.Numeric{Int: big.NewInt(23), Exp: 1, NaN: false, Valid: true}

	id, err := queries.InsertApplianceAndReturnId(c.Request().Context(), out.InsertApplianceAndReturnIdParams{
		UserID:  pgtype.UUID{Bytes: uu, Valid: true},
		Name:    "Computer",
		Wattage: wattage,
	})

	return id, err
}

func UpdateApplianceGroupID(c echo.Context, queries *out.Queries) {
	applianceId, err := util.ConvertStringTo16Bytes("c3018885a487467cb1f794dc72850ab1")
	if err != nil {
		c.Logger().Fatal("ERROR CONVERTING STRING TO 16 BYTE SLICE")
	}

	groupID, err := util.ConvertStringTo16Bytes("c3018885a487467cb1f794dc72850ab1")
	if err != nil {
		c.Logger().Fatal("ERROR CONVERTING STRING TO 16 BYTE SLICE")
	}

	err = queries.UpdateApplianceGroupID(c.Request().Context(), out.UpdateApplianceGroupIDParams{
		GroupID: pgtype.UUID{Bytes: groupID, Valid: true},
		ID:      pgtype.UUID{Bytes: applianceId, Valid: true},
	})

	if err != nil {
		c.Logger().Fatal("ERROR UPDATING GROUP")
	}
}
