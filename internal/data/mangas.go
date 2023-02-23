package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type MangaModel struct {
	DB *pgxpool.Pool
}

func (m MangaModel) Insert(manga *Manga) error {
	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := conn.QueryRow(ctx, "INSERT INTO mangas (title, studio, year, chapters, rating) VALUES ($1, $2, $3, $4, $5) RETURNING id, version", manga.Title, manga.Studio, manga.Year, manga.Chapters, manga.Rating)
	err = row.Scan(&manga.ID, &manga.Version)
	if err != nil {
		return err
	}

	return nil
}

func (m MangaModel) Get(id int64) (*Manga, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	var manga Manga

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := conn.QueryRow(ctx, "SELECT * FROM mangas WHERE id = $1", id)

	err = row.Scan(&manga.ID, &manga.Title, &manga.Studio, &manga.Year, &manga.Chapters, &manga.Rating, &manga.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &manga, nil
}

func (m MangaModel) Update(manga *Manga) error {
	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := conn.QueryRow(ctx, "UPDATE mangas SET title = $1, studio = $2, year = $3, chapters = $4, rating = $5, version = version + 1 WHERE id = $6 AND version = $7 RETURNING version", manga.Title, manga.Studio, manga.Year, manga.Chapters, manga.Rating, manga.ID, manga.Version)
	err = row.Scan(&manga.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m MangaModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := conn.Query(ctx, "DELETE FROM mangas WHERE id = $1", id)
	if err != nil {
		return err
	}
	err = res.Err()
	if err != nil {
		return err
	}
	return nil
}

func (m MangaModel) GetAll(title string, genres []string, filters Filters) ([]*Manga, Metadata, error) {
	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return nil, Metadata{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := conn.Query(ctx, "SELECT id, title, studio, year, chapters, rating, version FROM mangas")
	if err != nil {
		return nil, Metadata{}, err
	}

	defer res.Close()

	mangas := []*Manga{}

	for res.Next() {
		var manga Manga
		err := res.Scan(&manga.ID, &manga.Title, &manga.Studio, &manga.Year, &manga.Chapters, &manga.Rating, &manga.Version)
		if err != nil {
			return nil, Metadata{}, err
		}
		mangas = append(mangas, &manga)
	}

	if err = res.Err(); err != nil {
		return nil, Metadata{}, err
	}

	//metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return mangas, Metadata{}, nil
}

type Manga struct {
	ID       int64   `json:"id"`
	Title    string  `json:"title"`
	Studio   string  `json:"studio"`
	Year     int32   `json:"year,omitempty"`
	Chapters int32   `json:"chapters"`
	Rating   float32 `json:"rating"`
	Version  int32   `json:"version"`
}

/*
func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
*/
