package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"
	"strings"

	"snippetbox.glebich/internal/models"
	"snippetbox.glebich/internal/validator"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

type userSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.html", data)
	/*
		for _, snippet := range snippets {
			fmt.Fprintf(w, "%+v\n", snippet)
		}
		w.Write([]byte("Hello Go dev!"))
	*/
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// чтобы убрать экранированные знаки переноса строки
	snippet.Content = strings.Replace(snippet.Content, "\\n", "\n", -1)

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.html", data)

	//fmt.Fprintf(w, "Display a specific snippet with ID %d...\n", id)
	//fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	/*
		if r.Method == http.MethodOptions {
			w.Header().Set("Allow", "POST")
			return
		} else if r.Method != http.MethodPost {
			//w.WriteHeader(http.StatusMethodNotAllowed)
			//w.Write([]byte("POST only"))
			//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}
	*/
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignupGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := userSignupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}
	form.CheckField(validator.MaxChars(form.Name, 20), "name", "This field cannot be more than 20 characters long")
	form.CheckField(validator.ValidName(form.Name), "name", "This field can only contain letters, numbers and symbols - _")
	form.CheckField(validator.MinChars(form.Name, 3), "name", "This field must contain more than 3 characters")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Please enter correct email")
	form.CheckField(validator.ValidPassword(form.Password), "password", "Password must contain 1 number (0-9), 1 uppercase letter, 1 lowercase letter, 1 non-alpha numeric number, password is 8-16 characters with no space")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	id, err := app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = CreateJWTTokenAndSetCookie(form.Name, form.Email, id, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.GenerateRefreshTokenAndCookie(w, id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLoginGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	user, err := app.users.Get(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrWrongCredentials) {
			form.AddFieldError("credentials", "Wrong Credentials")
		} else {
			app.serverError(w, err)
			return
		}
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	err = CreateJWTTokenAndSetCookie(user.Name, user.Email, user.ID, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.GenerateRefreshTokenAndCookie(w, user.ID)
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := app.refreshTokens.Delete(data.User.ID)
	if err != nil {
		app.serverError(w, fmt.Errorf("fr I don't know WTF"))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // если HTTPS - true, локальная разработка - false
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // если HTTPS - true, локальная разработка - false
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	// наверно, можно передавать в контекст ещё разные сообщения,
	// чтобы информационные уведомления показывать типа
	// вы вышли из аккаунта, вы вошли в акк и тп
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
