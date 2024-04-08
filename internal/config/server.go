package config

import (
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/frcrod/watt-companion-api/api/handler"
	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/frcrod/watt-companion-api/internal/types"
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
	e.Use(middleware.RemoveTrailingSlash())

	key := os.Getenv("SESSION_SECRET")
	maxAge := 86400 * 1                      // 1 day
	isProd := os.Getenv("IS_PROD") == "true" // Set to true when serving over https

	gob.Register(types.AuthUser{})
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
	e.GET("/user", func(e echo.Context) error {
		session, err := gothic.Store.Get(e.Request(), "auth")
		if err != nil {
			slog.Error("Error getting the auth session")
			return err
		}

		marshalled, err := json.Marshal(session.Values["user"])
		if err != nil {
			slog.Error("Error marshalling the user session")
			return err
		}

		user := goth.User{}

		if err := json.Unmarshal(marshalled, &user); err != nil {
			slog.Error("Error reading the Marshalled Session")
			return err
		}

		fmt.Printf("marshalled: %v\n", user.Email)
		return e.String(http.StatusOK, string(marshalled))
	})

	e.GET("/auth/:provider/callback", func(e echo.Context) error { return handler.AuthCallback(e, queries) })
	e.GET("/logout/:provider", handler.Logout)
	e.GET("/auth/:provider", func(e echo.Context) error { return handler.DoAuth(e, queries) })

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

	e.GET("/appliance", handler.Home)
	e.Static("static", "web/static")

	return e
}
