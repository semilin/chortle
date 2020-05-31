package main

import (
	"github.com/mholt/archiver"
	"github.com/theckman/yacspin"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

// DownloadSong downloads input song
// to the user's song directory
func DownloadSong(song Song) {
	// generate directory name with format
	// Title - Artist (Charter)
	dirname := song.Name + " - " + song.Artist + " (" + song.Charter + ")"
	// create directory in
	// the songs folder
	err := os.Mkdir(path.Join(songsDir(), dirname), 0755)
	if HandleErr(err) == "fileexists" {
		os.RemoveAll(path.Join(songsDir(), dirname))
		os.Mkdir(path.Join(songsDir(), dirname), 0755)
	}

	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)
	HandleErr(err)

	fields := reflect.TypeOf(song.Downloads)
	vals := reflect.ValueOf(song.Downloads)

	spinner.Start()

	for i := 0; i < 11; i++ {
		field := fields.Field(i)  // get names and values
		val := vals.Field(i) 	  // of song downloads

		if val.String() != "" {
			// if field is filled, download the file
			spinner.Suffix(field.Name)
			spinner.Message("downloading")

			resp, err := http.Get(val.String())
			HandleErr(err)
			body, err := ioutil.ReadAll(resp.Body)
			HandleErr(err)
			defer resp.Body.Close()

			if strings.Contains(string(body), "div id=") {
				Message("There was an error with the download. Please try again later.")
			}

			var fn string // filename

			switch field.Name {
			case "INI":
				fn = "song.ini"
			case "Chart":
				fn = "notes.chart"
			case "Video":
				fn = "video.mp4"
			case "SongMP3":
				fn = "song.mp3"
			case "SongOGG":
				fn = "song.ogg"
			case "AlbumPNG":
				fn = "album.png"
			case "AlbumJPG":
				fn = "album.jpg"
			case "Drums":
				fn = "drums.ogg"
			case "GuitarOGG":
				fn = "guitar.ogg"
			case "GuitarMP3":
				fn = "guitar.mp3"
			case "Archive":
				mime := http.DetectContentType(body)
				var ext string

				switch mime {
				case "application/x-rar-compressed":
					ext = ".rar"
				case "application/zip":
					ext = ".zip"
				}
				fn = strings.Join([]string{"archive", ext}, "")

			}

			spinner.Message("writing")

			// Write the file from RAM
			// to the song folder

			writepath := path.Join(songsDir(), dirname, fn)
			ioutil.WriteFile(writepath, body, 0644)

			if strings.HasPrefix(fn, "archive") {
				// if it's an archive file,
				// extract it
				spinner.Message("extracting")

				writepath := path.Join(songsDir(), dirname, fn)

				err := archiver.Unarchive(writepath, path.Join(songsDir(), dirname))
				HandleErr(err)

				spinner.Message("cleaning up")
				os.Remove(writepath)
			}

		}
	}

	spinner.Suffix("done")

	spinner.Stop()
}

func songsDir() string {
	file, err := os.OpenFile("./songsdir", os.O_WRONLY, 0600)
	HandleErr(err)
	defer file.Close()

	dir, err := ioutil.ReadAll(file)
	HandleErr(err)
	return strings.TrimSpace(string(dir))
}
