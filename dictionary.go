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
	Results []Results `json:"results"`
}

type Senses struct {
	Definitions []string `json:"definitions"`
}

type Entries struct {
	Senses []Senses `json:"senses"`
}

type LexicalEntries struct {
	Entries []Entries `json:"entries"`
}

type Results struct {
	LexicalEntries []LexicalEntries `json:"lexicalEntries"`
}

func getDefinition(word string) string {
	appKey := os.Getenv("appKey")
	appId := os.Getenv("appId")
	baseUrl := "https://od-api.oxforddictionaries.com/api/v2"
	lang := "en-gb"

	fullUrl := fmt.Sprintf("%s/entries/%s/%s", baseUrl, lang, strings.ToLower(word))

	result := &oxford{}
	client := &http.Client{}

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Println("Get request to dictionary failed: ", err)
		return ""
	}
	req.Header.Set("app_id", appId)
	req.Header.Set("app_key", appKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Response issue: ", err)
		return "No definition Found"
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		log.Println("Json Decoder: ", err)
		return "Something is wrong on our end. Check back in a bit."
	}

	definition := "Definitions:"

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
	urbanUrl := "https://www.urbandictionary.com/define.php?term=" + word
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("x-rapidapi-key", os.Getenv("x-rapidapi-key"))
	req.Header.Set("x-rapidapi-host", os.Getenv("x-rapidapi-host"))

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "No definition found."
	}

	urbanResponse := &urban{}
	definition := "Definitions:"
	err = json.NewDecoder(resp.Body).Decode(urbanResponse)
	if err != nil {
		return "Something is wrong on our end. Check back in a bit."
	}

	i := 0
	for i < 6 {
		definition += "\n" + strconv.Itoa(i+1) + ". " + urbanResponse.List[i].Definition
		i ++
	}

	return definition + "\n\nCheck " + urbanUrl + "for more."
}
