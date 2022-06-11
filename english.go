package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

/* this uses https://dictionary.dev
which should always be free */

type Response struct {
	Word string `json:"word"`
	//Phonetic   string      `json:"phonetic"`
	//Phonetics  []Phonetics `json:"phonetics"`
	Meanings []Meanings `json:"meanings"`
	//License    License     `json:"license"`
	//SourceUrls []string    `json:"sourceUrls"`
}

/*type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Phonetics struct {
	Text      string  `json:"text"`
	Audio     string  `json:"audio"`
	SourceURL string  `json:"sourceUrl,omitempty"`
	License   License `json:"license,omitempty"`
}*/

type Definitions struct {
	Definition string `json:"definition"`
	//Synonyms   []string `json:"synonyms"`
	//Antonyms   []string `json:"antonyms"`
	//Example    string        `json:"example,omitempty"`
}
type Meanings struct {
	PartOfSpeech string        `json:"partOfSpeech"`
	Definitions  []Definitions `json:"definitions"`
	//Synonyms     []string      `json:"synonyms"`
	//Antonyms     []string `json:"antonyms"`
}

func getMeaning(word string) string {
	const baseUrl = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	fullUrl := baseUrl + word

	var response []Response
	client := &http.Client{}

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Println(err)
		return ""
	}

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println(err)
		suggest := grammarChecker(word, 5)
		return suggest
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return "Something went wrong on our end. Check back in a bit."
	}

	result := "<b>Entry:</b> " + word

	entry := response[0]
	noOfMeanings := 10 / len(entry.Meanings)

	for i := 0; i < len(entry.Meanings); i++ {
		goingTo := noOfMeanings
		if len(entry.Meanings[i].Definitions) < goingTo {
			goingTo = len(entry.Meanings[i].Definitions)
		}

		result += "\n\n<b>Part of Speech: </b>" + entry.Meanings[i].PartOfSpeech
		for j := 0; j < goingTo; j++ {
			result += "\n" + strconv.Itoa(j+1) + ". " + entry.Meanings[i].Definitions[j].Definition
		}
	}

	return result
}
