package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Sqlite3 Test")
	db := NewDb("./MyVideos107.db")
	// fmt.Printf("%v", db.GetMovie())

	// DebugMovieArray(db.GetMovie())
	// DebugMovieStruct(db.GetMovie()[1])

	ans := Input("Search Movies : ")
	// fmt.Println(ans)
	DebugMovieArray(db.GetMovieWithQuery(strings.Join(ans, " ")))
}
