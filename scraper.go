package main

import (
	"github.com/theckman/yacspin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Search takes any term and looks it up in
// the Chorus database.
func Search(q string) []byte {
	chorus := "http://chorus.fightthe.pw/api/search?query="
	query := parseQuery(q)

	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          "searching",
		SuffixAutoColon: true,
		Message:         "",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)
	HandleErr(err)

	spinner.Start()

	raw := getPage(chorus + query)
	spinner.Stop()
	return raw
}

// Get rid of spaces to make a query readable
// by the API.
func parseQuery(q string) string {
	q = strings.ReplaceAll(q, " ", "%20")
	return q
}

func getPage(url string) []byte {
	resp, err := http.Get(url)
	HandleErr(err)

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	HandleErr(err)

	return html
}
