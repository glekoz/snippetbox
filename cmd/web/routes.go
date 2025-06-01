package main

import (
	"net/http"

	"github.com/justinas/alice"
	"snippetbox.glebich/ui"
)

func (app *application) routes() http.Handler {
	// глобальный handler - распределяет запросы по другим обработчикам
	mux := http.NewServeMux()

	/*
		// для обработки статических файлов - картинок, css и т.д.
		fileServer := http.FileServer(http.Dir("./ui/static/"))
		mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	*/
	// обработчики конкретных путей
	// чтобы так на каждый эндпоинт не писать - сделать новую цепочку

	// для обработки статических файлов - картинок, css и т.д.
	fileserver := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileserver)

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

	// noSurf глобально, так как лог аут находится в нав баре, можно выйти из любой страницы,
	// так что нужно везде вставлять csrf токен в куки
	// может, только для файл сервера не надо
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders, app.authenticate, app.noSurf)

	return standard.Then(mux)
}
