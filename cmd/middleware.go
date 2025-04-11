package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/sangketkit01/personal-block/internal/helpers"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler{
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and save the session on every request
func SessionLoad(next http.Handler) http.Handler{
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !helpers.IsAuthenticated(r){
			session.Put(r.Context(), "error" , "Login first")
			http.Redirect(w,r,"/login",http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w,r)
	})
}

func PreventAuthGoSignUpAndSignIn(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if app.Session.Exists(r.Context(),"user"){
			http.Redirect(w,r,"/",http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w,r)
	})
}