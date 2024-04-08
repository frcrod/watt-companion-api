package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/frcrod/watt-companion-api/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func addUserToSession(wr http.ResponseWriter, req *http.Request, user types.AuthUser) {
	session, err := gothic.Store.Get(req, "auth")
	if err != nil {
		log.Print("Error ", err)
	}

	// Remove the raw data to reduce the size
	user.User.RawData = map[string]interface{}{}

	session.Values["user"] = user

	err = session.Save(req, wr)
	if err != nil {
		log.Print("Problem Saving session data", err)
	}
}

func DoAuth(e echo.Context, queries *out.Queries) error {
	provider := e.Param("provider")

	r := e.Request().WithContext(context.WithValue(context.Background(), "provider", provider))
	if gothUser, err := gothic.CompleteUserAuth(e.Response().Writer, r); err == nil {
		session, err := gothic.Store.Get(e.Request(), "auth")

		if err != nil {
			log.Fatal("Error getting session")
			return err
		}

		session.Values["user"] = gothUser
		err = session.Save(e.Request(), e.Response().Writer)
		if err != nil {
			log.Fatal("Error saving the user session")
		}

		return e.String(200, gothUser.Name)
	} else {
		gothic.BeginAuthHandler(e.Response().Writer, r)
	}

	return e.String(200, "AUTH")
}

func Logout(e echo.Context) error {
	session, err := gothic.Store.Get(e.Request(), "auth")
	if err != nil {
		log.Fatal("Error getting session")
		return err
	}

	session.Values["user"] = goth.User{}
	session.Options.MaxAge = -1
	session.Save(e.Request(), e.Response().Writer)

	e.Response().Header().Set("Location", "/")
	e.Response().WriteHeader(http.StatusTemporaryRedirect)

	return nil
}

func AuthCallback(e echo.Context, queries *out.Queries) error {
	provider := e.Param("provider")

	r := e.Request().WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(e.Response().Writer, r)
	if err != nil {
		log.Fatal("Error logging in", err)
		return err
	}

	id, err := queries.CheckUserExists(e.Request().Context(), user.Email)

	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			id, err = queries.CreateUserAndReturnId(r.Context(), out.CreateUserAndReturnIdParams{
				Email:    user.Email,
				Nickname: user.FirstName,
			})

			if err != nil {
				log.Fatal("Error creating user", err)
				return err
			}
		}
	}

	uu, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		log.Fatal("Error parsing uuid", err)
	}

	addUserToSession(e.Response().Writer, e.Request(), types.AuthUser{
		Id:              uu,
		User:            user,
		IsAuthenticated: true,
	})
	e.Response().Header().Set("Location", "/")
	e.Response().WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
