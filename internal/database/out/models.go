// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package out

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Appliance struct {
	ID        pgtype.UUID
	Name      string
	Wattage   pgtype.Numeric
	GroupID   pgtype.UUID
	UserID    pgtype.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Group struct {
	ID        pgtype.UUID
	Name      string
	UserID    pgtype.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type User struct {
	ID        pgtype.UUID
	Username  string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}