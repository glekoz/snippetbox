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
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)

	protected := alice.New(app.requireAuth)
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreateGet))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))

	altProtected := alice.New(app.requireNoAuth)
	mux.Handle("GET /user/signup", altProtected.ThenFunc(app.userSignupGet))
	mux.Handle("POST /user/signup", altProtected.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", altProtected.ThenFunc(app.userLoginGet))
	mux.Handle("POST /user/login", altProtected.ThenFunc(app.userLoginPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.authenticate)

	return standard.Then(mux)
}
