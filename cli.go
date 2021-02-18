package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	args := os.Args
	var query string
	if len(args) == 1 {
		query = PromptSearchSong()
	} else {
		query = strings.Join(args[1:], " ")
	}
	raw := Search(query)
	songs := ToSongs(raw)
	i := PromptDownloadSong(songs)
	DownloadSong(songs.Songs[i])
	
}

func PromptFileOverwrite() bool {
	resp := prompt("You have this chart. Do you want to overwrite it? (y/n)")
	if resp == "y" {
		return true
	}
	return false
}

func PromptDownloadSong(songs Songs) int {
	var data [][]string
	for i := range songs.Songs {
		song := songs.Songs[i]
		index := strconv.Itoa(i+1)
		var nrow []string
		nrow = append(nrow, index)
		nrow = append(nrow, song.Name)
		nrow = append(nrow, song.Artist)
		nrow = append(nrow, song.Charter)

		data = append(data, nrow)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Name", "Artist", "Charter"})
	table.AppendBulk(data)
	table.Render()

	resp := prompt("Pick a song to download (1-20) (default 1)")
	var iresp int
	if resp != "" {
		var err error
		iresp, err = strconv.Atoi(resp)
		HandleErr(err)
		iresp-- 
	} else {
		iresp = 0
	}

	return iresp
}

func PromptSearchSong() string {
	resp := prompt("Enter a search term")
	return resp
}

func PromptSongDir() {
	resp := prompt("Please enter your songs directory")
	err := ioutil.WriteFile("./songsdir", []byte(resp), 0644)
	HandleErr(err)
}

func Message(s string) {
	fmt.Println(s)
}

func HandleErr(err error) string {
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			if PromptFileOverwrite() {
				return "file_exists_overwrite"
			} else {
				return "file_exists_keep"
			}
		} else if os.IsNotExist(err) {
			PromptSongDir()
		} else {
			panic(err)
		}
	}
	return ""
}

func prompt(s string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", s)
	resp, err := reader.ReadString('\n')
	HandleErr(err)
	clean := strings.TrimSpace(resp)
	return clean
}
