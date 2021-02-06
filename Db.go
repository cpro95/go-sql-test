package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// Db type for sqlite3 DB Object
type Db struct {
	db *sql.DB
}

// NewDb returns db which name is dbName
func NewDb(dbName string) *Db {
	// Openning db file
	database, err := sql.Open("sqlite3", dbName)

	// Error
	if err != nil {
		log.Fatal(err)
	}

	d := &Db{
		db: database,
	}

	log.Info("DB OPENED")
	return d
}

// Close is for db close
func (d *Db) Close() {
	log.Info("DB CLOSED")
	d.db.Close()
}

// GetMovie return Movies array
func (d *Db) GetMovie() []Movie {
	var movie Movie
	var movies []Movie
	var err error

	rows, err := d.db.Query("select idMovie, c00, c01, c03, premiered, strPath, rating from movie_view where c00 like '%star%'")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&movie.IDMovie, &movie.C00, &movie.C01, &movie.C03, &movie.Premiered, &movie.StrPath, &movie.Rating)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	rows.Close()

	return movies
}

// GetMovieWithQuery return Movies with query search result
func (d *Db) GetMovieWithQuery(query string) []Movie {
	var movie Movie
	var movies []Movie
	var err error

	query = "select idMovie, c00, c01, c03, premiered, strPath, rating from movie_view where c00 LIKE '%" + query + "%' order by idMovie"
	rows, err := d.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&movie.IDMovie, &movie.C00, &movie.C01, &movie.C03, &movie.Premiered, &movie.StrPath, &movie.Rating)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	rows.Close()

	return movies
}
