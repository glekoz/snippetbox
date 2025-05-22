package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// глобальный handler - распределяет запросы по другим обработчикам
	mux := http.NewServeMux()

	// для обработки статических файлов - картинок, css и т.д.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// обработчики конкретных путей
	mux.HandleFunc("GET /", app.home)
	mux.Handle("GET /snippet/view/{id}", app.requestJWT(http.HandlerFunc(app.snippetView)))
	mux.HandleFunc("GET /snippet/create", app.snippetCreateGet)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)
	mux.HandleFunc("GET /user/signup", app.userSignupGet)
	mux.HandleFunc("POST /user/signup", app.userSignupPost)
	mux.HandleFunc("GET /user/login", app.userLoginGet)
	mux.HandleFunc("POST /user/login", app.userLoginPost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
