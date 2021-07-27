package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type oxford struct {
	Results []results `json:"results"`
}

type senses struct {
	Definitions []string `json:"definitions"`
}

type entries struct {
	Senses []senses `json:"senses"`
}

type lexicalEntries struct {
	Entries []entries `json:"entries"`
}

type results struct {
	LexicalEntries []lexicalEntries `json:"lexicalEntries"`
}

func getDefinition(word string) string {
	appKey := os.Getenv("appKey")
	appID := os.Getenv("appId")
	baseURL := "https://od-api.oxforddictionaries.com/api/v2"
	lang := "en-gb"

	fullURL := fmt.Sprintf("%s/entries/%s/%s", baseURL, lang, strings.ToLower(word))

	result := &oxford{}
	client := &http.Client{}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Println("Get request to dictionary failed: ", err)
		return ""
	}
	req.Header.Set("app_id", appID)
	req.Header.Set("app_key", appKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Response issue: ", err)
		return "<b>Entry:</b> " + strings.ToTitle(word) + "\nNo definition Found."
	}

	if resp.StatusCode == http.StatusNotFound {
		suggest := grammarChecker(word, 5)
		return suggest
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		log.Println("Json Decoder: ", err)
		return "Something is wrong on our end. Check back in a bit."
	}

	definition := "<b>Entry:</b> " + strings.ToTitle(word) +"\n\n<b>Definition(s):</b>"

	count := 1

	for _, result := range result.Results {
		for _, entry := range result.LexicalEntries {
			for _, senses := range entry.Entries {
				for _, sens := range senses.Senses {
					for _, def := range sens.Definitions {
						definition += "\n" + strconv.Itoa(count) + ". " + strings.Title(def)
						count++
					}
				}
			}
		}
	}

	return definition
}

// urban dictionary
type urban struct {
	List	[]list	`json:"list"`
}

type list struct {
	Definition string	`json:"definition"`
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
		return "<b>Entry:</b> " + strings.ToTitle(word) + "\n\nNo definition found."
	}

	urbanResponse := &urban{}
	definition := "<b>Entry:</b> " + strings.ToTitle(word) +"\n\n<b>Definition(s):</b>"
	err = json.NewDecoder(resp.Body).Decode(urbanResponse)
	if err != nil {
		return "Something is wrong on our end. Check back in a bit."
	}

	i := 0
	max := 5

	if len(urbanResponse.List) == 0 {
		return "<b>Entry:</b> " + strings.ToTitle(word) + "\nNo definitions found."
	}
	if len(urbanResponse.List) < 5 {
		max = len(urbanResponse.List)
	}
	for i < max {
		definition += "\n" + strconv.Itoa(i+1) + ". " + urbanResponse.List[i].Definition
		i ++
	}

	return definition + "\n\nCheck " + urbanURL + " for more."
}
