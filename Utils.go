package main

import "github.com/davecgh/go-spew/spew"

// DebugMovieArray dump array struct
func DebugMovieArray(arr []Movie) {
	for _, element := range arr {
		spew.Dump(element)
	}
}

// DebugMovieStruct dump struct
func DebugMovieStruct(str Movie) {
	spew.Dump(str)
}
