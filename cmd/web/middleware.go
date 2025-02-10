package main

import (
	"net/http"
)

func (app *Application) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;")
		next.ServeHTTP(w, r)
	})
}

func (app *Application) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r.Header.Get("token"), app.secretKeyJWT) {
			setDescriptionStatusCode("The user is not logged in", http.StatusUnauthorized, w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := VerifyToken(r.Header.Get("token"), app.secretKeyJWT)
		if err != nil || role != "admin" {
			setDescriptionStatusCode("The user does not have access", http.StatusForbidden, w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
