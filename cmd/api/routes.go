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

	//router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	//router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	//router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	// router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	// router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
	router.HandlerFunc(http.MethodPost, "/v1/mangas", app.createMangaHandler)
	router.HandlerFunc(http.MethodGet, "/v1/mangas", app.listMangaHandler)
	router.HandlerFunc(http.MethodGet, "/v1/mangas/:id", app.showMangaHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/mangas/:id", app.updateMangaHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/mangas/:id", app.deleteMangaHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/", app.registerUserHandler)

	//	return app.recoverPanic(app.rateLimit(app.autenticate(router)))
	return app.recoverPanic(app.rateLimit(router))
}
