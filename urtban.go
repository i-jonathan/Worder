package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// urban dictionary
type urban struct {
	List []list `json:"list"`
}

type list struct {
	Definition string `json:"definition"`
}

func getUrbanDefinition(word string) string {
	url := "https://mashape-community-urban-dictionary.p.rapidapi.com/define?term=" + word
	urbanURL := "https://www.urbandictionary.com/define.php?term=" + word
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("x-rapidapi-key", os.Getenv("x-rapidapi-key"))
	req.Header.Set("x-rapidapi-host", os.Getenv("x-rapidapi-host"))

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "<b>Entry:</b> " + strings.Title(word) + "\n\nNo definition found."
	}

	urbanResponse := &urban{}
	definition := "<b>Entry:</b> " + strings.Title(word) + "\n\n<b>Definition(s):</b>"
	err = json.NewDecoder(resp.Body).Decode(urbanResponse)
	if err != nil {
		return "Something is wrong on our end. Check back in a bit."
	}

	i := 0
	max := 5

	if len(urbanResponse.List) == 0 {
		return "<b>Entry:</b> " + strings.Title(word) + "\nNo definitions found."
	}
	if len(urbanResponse.List) < 5 {
		max = len(urbanResponse.List)
	}
	for i < max {
		definition += "\n" + strconv.Itoa(i+1) + ". " + urbanResponse.List[i].Definition
		i++
	}

	return definition + "\n\nCheck " + urbanURL + " for more."
}
