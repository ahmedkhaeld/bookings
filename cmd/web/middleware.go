package main

import (
	"github.com/ahmedkhaeld/bookings/internal/helpers"
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf adds CSRF protection against all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request for a middleware
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth a middleware to protect the routes only for logged in
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
