package main

import "fmt"

// Movie type for json struct
type Movie struct {
	IDMovie   int         `json:"id"`
	C00       string      `json:"title"`
	C01       string      `json:"overview"`
	C03       string      `json:"tagline"`
	Premiered string      `json:"premiered"`
	StrPath   string      `json:"strPath"`
	Rating    interface{} `json:"rating"`
}

// NewMovie returns Movie Struct
func NewMovie() *Movie {
	m := &Movie{}
	return m
}

// ToString is return Movie Struct
func (m *Movie) ToString() string {
	return fmt.Sprintf("%v", m)
}
