package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	standardChain := alice.New(app.LogRequest, app.JsonContentTypeMiddleware, app.RequireAuthentication)
	admin := standardChain.Append(app.AdminOnly)

	router.HandlerFunc(http.MethodGet, "/user_banner", app.getUserBanner)

	router.Handler(http.MethodGet, "/banner", admin.ThenFunc(app.getBanner))
	router.Handler(http.MethodPost, "/banner", admin.ThenFunc(app.postBanner))
	router.Handler(http.MethodPatch, "/banner/:id", admin.ThenFunc(app.patchBannerByID))
	router.Handler(http.MethodDelete, "/banner/:id", admin.ThenFunc(app.deleteBannerByID))
	return standardChain.Then(router)
}
