package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/mangas", app.requirePermission("movies:write", app.createMangaHandler))
	router.HandlerFunc(http.MethodGet, "/v1/mangas", app.requirePermission("movies:read", app.listMangaHandler))
	router.HandlerFunc(http.MethodGet, "/v1/mangas/:id", app.requirePermission("movies:read", app.showMangaHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/mangas/:id", app.requirePermission("movies:write", app.updateMangaHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/mangas/:id", app.requirePermission("movies:write", app.deleteMangaHandler))

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(app.autenticate(router)))
}
