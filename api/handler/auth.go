package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func DoAuth(e echo.Context) error {
	provider := e.Param("provider")

	r := e.Request().WithContext(context.WithValue(context.Background(), "provider", provider))
	if gothUser, err := gothic.CompleteUserAuth(e.Response().Writer, r); err == nil {
		session, err := gothic.Store.Get(e.Request(), "auth")

		if err != nil {
			log.Fatal("Error getting session")
			return err
		}

		fmt.Println(gothic.GetFromSession(gothUser.UserID, e.Request()))
		session.Values["user"] = gothUser
		err = session.Save(e.Request(), e.Response().Writer)
		if err != nil {
			log.Fatal("Error saving the user session")
		}

		fmt.Println(session.Store().Get(e.Request(), "user"))
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

func AuthCallback(e echo.Context) error {
	provider := e.Param("provider")

	r := e.Request().WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(e.Response().Writer, r)
	if err != nil {
		fmt.Fprintln(e.Response(), err)
		return err
	}

	return e.String(200, user.Email)
}
