package main

import (
	"encoding/json"
)

// Downloads is a simple wrapper type for
// storing download links for a song
type Downloads struct {
	INI       string `json:"ini"`
	Chart     string `json:"chart"`
	Video     string `json:"video"`
	SongOGG   string `json:"song.ogg"`
	SongMP3   string `json:"song.mp3"`
	AlbumPNG  string `json:"album.png"`
	AlbumJPG  string `json:"album.jpg"`
	Drums     string `json:"drums"`
	GuitarOGG string `json:"guitar.ogg"`
	GuitarMP3 string `json:"guitar.mp3"`
	Archive   string `json:"archive"`
}

// Song is a datatype that is much more easily
// parsed than the pure JSON.
// By creating these structures beforehand, it
// makes writing code that needs to get certain
// elements much more pleasant.
type Song struct {
	Name      string    `json:"name"`
	Artist    string    `json:"artist"`
	Album     string    `json:"album"`
	Genre     string    `json:"genre"`
	Charter   string    `json:"charter"`
	Downloads Downloads `json:"directLinks"`
}

// Songs is a struct that only contains an
// array of Song structs.
// This is quite odd, but I'm only doing it
// as to cooperate with the json module's
// very rigid Unmarshal function.
type Songs struct {
	Songs []Song `json:"songs"`
}

// ToSongs converts JSON to a Songs struct.
func ToSongs(b []byte) Songs {
	var songs Songs
	json.Unmarshal(b, &songs)

	return songs
}
