package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func handler(resp http.ResponseWriter, req *http.Request) {
	body := &webHookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		log.Println("Could not decode the request body ", err)
		return
	}

	go processRequest(body)

}

// process request
func processRequest(update *webHookReqBody) {
	parts := strings.Fields(update.Message.Text)

	helpText := "Supported commands:\n/english word - Define word with British English Dictionary\n/urban word " +
		"- Define word with Urban Dictionary\n/help - Display this help text"

	// if len(parts) > 2 {
	// 	if err := respond(update.Message.Chat.ID, "Please check your message and resend"); err != nil {
	// 		log.Println("Error in sending message ", err)
	// 		return
	// 	}
	// }

	switch len(parts) {
	case 1:
		command := parts[0]
		switch command {
		case "/start":
			welcomeMessage := "Hi.\nWelcome to Worder.\n\n" + helpText
			if err := respond(update.Message.Chat.ID, welcomeMessage); err != nil {
				log.Println("Error in sending message ", err)
				return
			}
		case "/help":
			if err := respond(update.Message.Chat.ID, helpText); err != nil {
				log.Println("Error in sending message ", err)
				return
			}
		default:
			if err := respond(update.Message.Chat.ID, helpText); err != nil {
				log.Println("Error in sending message ", err)
				return
			}
		}
	case 2:
		command := parts[0]
		word := parts[1]
		switch command {
		case "/urban":
			definition := getUrbanDefinition(word)
			if err := respond(update.Message.Chat.ID, definition); err != nil {
				log.Println("Error in sending message ", err)
				return
			}
		case "/english":
			definition := getDefinition(word)
			if err := respond(update.Message.Chat.ID, definition); err != nil {
				log.Println("Error in sending message ", err)
				return
			}
		default:
			if err := respond(update.Message.Chat.ID, helpText); err != nil {
				log.Println("Error in sending message ", err)
				return
			}
		}
	default:
		// if err := respond(update.Message.Chat.ID, helpText); err != nil {
		// 	log.Println("Error in sending message ", err)
		// 	return
		// }
		command := parts[0]
		word := ""

		for i := 1; i < len(parts); i++ {
			word += parts[i]
		}
		switch command {
			case "/urban":
				definition := getUrbanDefinition(word)
				if err := respond(update.Message.Chat.ID, definition); err != nil {
					log.Println("Error in sending message ", err)
					return
				}
			default:
				if err := respond(update.Message.Chat.ID, helpText); err != nil {
					log.Println("Error in sending message ", err)
					return
				}
		}
	}

	log.Println("Response sent")
}

// respond
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

type grammarSuggestions struct {
	Suggestions    []string `json:"suggestions"`
}

func grammarChecker(word string, entryCount int) string {
	url := fmt.Sprintf("https://api.collinsdictionary.com/api/v1/dictionaries/english/search/didyoumean?q=%s&entrynumber=%d", word, entryCount)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "No definitions or suggestions found."
	}

	req.Header.Set("accessKey", os.Getenv("accessKey"))

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "No definitions or suggestions found."
	}

	suggested := &grammarSuggestions{}

	err = json.NewDecoder(resp.Body).Decode(suggested)

	if err != nil {
		return "Something is wrong on our end. Give us a few."
	}

	words := "No definition found.\nDid you mean any of these words?\n"
	count := 1

	for _, suggestion := range suggested.Suggestions {
		words += fmt.Sprintf("%d. %s\n", count, suggestion)
		count ++
	}

	return words
}