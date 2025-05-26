package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"snippetbox.glebich/internal/jwtAuth"
)

type contextKey string

const contextKeyUser = contextKey("user")

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		var user *jwtAuth.Sub

		token, err := r.Cookie("auth_token")
		if err != nil {
			user, err = app.VerifyRefreshTokenAndCreateJWT(refreshToken.Value, w)
			if err != nil {
				if errors.Is(err, jwtAuth.ErrInvalidRefreshToken) {
					next.ServeHTTP(w, r)
					return
				}
				if errors.Is(err, jwtAuth.ErrServerError) {
					next.ServeHTTP(w, r)
					return
				} else {
					next.ServeHTTP(w, r)
					return
				}
			}
		} else {
			user, err = jwtAuth.VerifyJWTToken(token.Value)
			if err != nil {
				user, err = app.VerifyRefreshTokenAndCreateJWT(refreshToken.Value, w)
				if err != nil {
					next.ServeHTTP(w, r) // подробные ошибки добавить
					return
				}
			}
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(contextKeyUser).(*jwtAuth.Sub)
		if !ok {
			app.clientError(w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireNoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(contextKeyUser).(*jwtAuth.Sub)
		if ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
