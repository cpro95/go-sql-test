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
		log.Panic(err)
	}

	d := &Db{
		db: database,
	}

	return d
}

// Close is for db close
func (d *Db) Close() {
	d.db.Close()
}

// GetMovie is return Movies array
func (d *Db) GetMovie() []Movie {
	var movie Movie
	var movies []Movie
	var err error

	rows, err := d.db.Query("select idMovie, c00, c01, c03, premiered, rating from movie_view limit 2")
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		err := rows.Scan(&movie.IDMovie, &movie.C00, &movie.C01, &movie.C03, &movie.Premiered, &movie.Rating)
		if err != nil {
			log.Panic(err)
		}
		movies = append(movies, movie)
	}

	rows.Close()

	// log.Info(movies)

	return movies
}
