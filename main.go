package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type webHookReqBody struct {
	Message Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type reply struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

var token = os.Getenv("token")
var url = "https://api.telegram.org/bot" + token + "/"

func main() {
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, http.HandlerFunc(handler))
	log.Println(err)
}

func handler(res http.ResponseWriter, req *http.Request) {
	body := &webHookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		log.Println("Could not decode the request body ", err)
		return
	}

	if body.Message.Text == "/start" {
		welcomeMessage := "Hi.\nWelcome to Worder.\nTo get a definition, send an english word(en-gb) without unnecessary punctuations"
		if err := respond(body.Message.Chat.ID, welcomeMessage); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
	} else {
		definition := getDefinition(body.Message.Text)
		if err := respond(body.Message.Chat.ID, definition); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
	}

	log.Println("Response sent")
}

func respond(chatID int64, response string) error {
	reqBody := &reply{
		ChatID: chatID,
		Text:   response,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(url+"sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("Unexpected status " + res.Status)
	}

	return nil
}

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
		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		log.Println("Json Decoder: ", err)
		return ""
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
