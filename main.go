package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Sqlite3 Test")
	db := NewDb("./MyVideos107.db")
	// fmt.Printf("%v", db.GetMovie())

	DebugMovieArray(db.GetMovie())
	DebugMovieStruct(db.GetMovie()[1])
}
