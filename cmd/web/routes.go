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
	// чтобы так на каждый эндпоинт не писать - сделать новую цепочку
	mux.HandleFunc("GET /", app.home)
	mux.Handle("GET /snippet/view/{id}", alice.New(app.requireAuth).ThenFunc(app.snippetView))
	mux.HandleFunc("GET /snippet/create", app.snippetCreateGet)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)
	mux.Handle("GET /user/signup", alice.New(app.requireNoAuth).ThenFunc(app.userSignupGet))
	mux.Handle("POST /user/signup", alice.New(app.requireNoAuth).ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", alice.New(app.requireNoAuth).ThenFunc(app.userLoginGet))
	mux.Handle("POST /user/login", alice.New(app.requireNoAuth).ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", alice.New(app.requireAuth).ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.authenticate)

	return standard.Then(mux)
}
