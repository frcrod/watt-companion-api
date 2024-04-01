package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/frcrod/watt-companion-api/api/handler"
	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func CreateEchoInstance(queries *out.Queries) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CSRF())

	key := os.Getenv("SESSION_SECRET")
	maxAge := 86400 * 10                     // 30 days
	isProd := os.Getenv("IS_PROD") == "true" // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	googleClientId := os.Getenv("GOOGLE_KEY")
	googleClientSecret := os.Getenv("GOOGLE_SECRET")

	goth.UseProviders(google.New(googleClientId, googleClientSecret, "http://localhost:8080/auth/google/callback"))

	e.GET("/", handler.HandleLanding)
	e.POST("/user", func(c echo.Context) error {
		id, err := handler.AddUser(c, queries)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case pgerrcode.UniqueViolation:
					return c.String(400, "USERNAME EXISTS")
				}
			}
			return c.String(500, "ERROR")
		}

		return c.String(200, fmt.Sprintf("%s", hex.EncodeToString(id.Bytes[:])))
	})

	e.GET("/auth/:provider/callback", handler.AuthCallback)
	e.GET("/logout/:provider", handler.Logout)
	e.GET("/auth/:provider", handler.DoAuth)

	e.POST("/login", func(c echo.Context) error {
		id, err := handler.AddUser(c, queries)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case pgerrcode.UniqueViolation:
					return c.String(400, "USERNAME EXISTS")
				}
			}
			return c.String(500, "ERROR")
		}

		return c.String(200, fmt.Sprintf("%s", hex.EncodeToString(id.Bytes[:])))
	})

	e.POST("/appliance", func(c echo.Context) error {
		id, err := handler.AddAppliance(c, queries)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				}
			}
			return c.String(500, "ERROR")
		}

		return c.String(200, fmt.Sprintf("%s", hex.EncodeToString(id.Bytes[:])))
	})
	e.Static("static", "web/static")

	return e
}
