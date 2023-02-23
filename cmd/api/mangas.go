package main

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.adi.net/internal/validator"

	"greenlight.adi.net/internal/data"
)

func (app *application) createMangaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string  `json:"title"`
		Studio   string  `json:"studio"`
		Year     int32   `json:"year"`
		Chapters int32   `json:"chapters"`
		Rating   float32 `json:"rating"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Manga{
		Title:    input.Title,
		Studio:   input.Studio,
		Year:     input.Year,
		Chapters: input.Chapters,
		Rating:   input.Rating,
	}

	v := validator.New()

	if data.ValidateManga(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Mangas.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/mangas/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"manga": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showMangaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Mangas.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"manga": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateMangaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	manga, err := app.models.Mangas.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title    *string  `json:"title"`
		Studio   *string  `json:"studio"`
		Year     *int32   `json:"year"`
		Chapters *int32   `json:"chapters"`
		Rating   *float32 `json:"rating"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		manga.Title = *input.Title
	}
	if input.Studio != nil {
		manga.Studio = *input.Studio
	}
	if input.Year != nil {
		manga.Year = *input.Year
	}
	if input.Chapters != nil {
		manga.Chapters = *input.Chapters
	}
	if input.Rating != nil {
		manga.Rating = *input.Rating
	}
	v := validator.New()

	if data.ValidateManga(v, manga); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Mangas.Update(manga)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"manga": manga}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMangaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Mangas.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "manga successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMangaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	/*
		v := validator.New()
		qs := r.URL.Query()

		input.Title = app.readString(qs, "title", "")
		input.Genres = app.readCSV(qs, "genres", []string{})

		input.Filters.Page = app.readInt(qs, "page", 1, v)
		input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

		input.Filters.Sort = app.readString(qs, "sort", "id")
		input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

		if data.ValidateFilters(v, input.Filters); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}*/

	movies, metadata, err := app.models.Mangas.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"mangas": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
