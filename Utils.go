package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

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

// Input for input text
func Input(str string) []string {
	arr := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(str)
	scanner.Scan()
	text := scanner.Text()
	if len(text) != 0 {
		arr = append(arr, text)
	}
	return arr
}
