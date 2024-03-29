package main

/*
Discarding this cause oxford decided I can only check 1k words per eternity on a free plan.
Oh, well.

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
		return "<b>Entry:</b> " + strings.Title(word) + "\nNo definition Found."
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

	definition := "<b>Entry:</b> " + strings.Title(word) +"\n\n<b>Definition(s):</b>"

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
*/
