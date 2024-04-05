package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/frcrod/watt-companion-api/internal/database/out"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

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

	e.Response().Header().Set("Location", "/")
	e.Response().WriteHeader(http.StatusTemporaryRedirect)

	session.Values["user"] = goth.User{}
	session.Save(e.Request(), e.Response().Writer)

	return nil
}

func AuthCallback(e echo.Context, queries *out.Queries) error {
	provider := e.Param("provider")

	r := e.Request().WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(e.Response().Writer, r)
	if err != nil {
		log.Fatal("Error logging in")
		return err
	}

	_, err = queries.CheckUserExists(e.Request().Context(), user.Email)

	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			_, err = queries.CreateUserAndReturnId(r.Context(), out.CreateUserAndReturnIdParams{
				Email:    user.Email,
				Nickname: user.FirstName,
			})

			if err != nil {
				log.Fatal("Error creating user", err)
				return err
			}
		}
	}

	addUserToSession(e.Response().Writer, e.Request(), user)
	e.Response().Header().Set("Location", "/")
	e.Response().WriteHeader(http.StatusTemporaryRedirect)

	return nil
}

func addUserToSession(wr http.ResponseWriter, req *http.Request, user goth.User) {
	session, err := gothic.Store.Get(req, "auth")
	if err != nil {
		log.Print("Error ", err)
	}

	// Remove the raw data to reduce the size
	user.RawData = map[string]interface{}{}

	session.Values["user"] = user
	err = session.Save(req, wr)
	if err != nil {
		log.Print("Problem Saving session data", err)
	}
}
