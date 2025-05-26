package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"snippetbox.glebich/internal/jwtAuth"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// выводится вторая запись в трейсе, а не первая:
	// первая запись указывает на строку в этом файле, а вторая - откуда
	// была вызвана эта строка, то есть конкретное место ошибки
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	td := &templateData{
		CurrentYear: time.Now().Year(),
	}
	if user, ok := r.Context().Value(contextKeyUser).(*jwtAuth.Sub); ok {
		td.User = user
	}
	return td
}

func CreateJWTTokenAndSetCookie(name, email string, id int, w http.ResponseWriter) error {
	tokenString, err := jwtAuth.CreateJWTToken(name, email, id)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // если HTTPS - true, локальная разработка - false
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
	})

	return nil
}

func (app *application) VerifyRefreshTokenAndCreateJWT(refreshTokenString string, w http.ResponseWriter) (*jwtAuth.Sub, error) {
	user, err := app.refreshTokens.CheckRefreshToken(refreshTokenString)
	if err != nil {
		return nil, jwtAuth.ErrInvalidRefreshToken
	}
	err = CreateJWTTokenAndSetCookie(user.Name, user.Email, user.ID, w)
	if err != nil {
		return nil, jwtAuth.ErrServerError
	}
	return user, nil
}

func (app *application) GenerateRefreshTokenAndCookie(w http.ResponseWriter, userId int) error {
	refreshTokenString := rand.Text()

	err := app.refreshTokens.Insert(refreshTokenString, 1, userId)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // если HTTPS - true, локальная разработка - false
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
	})
	return nil
}
