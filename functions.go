package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
)

func handler(resp http.ResponseWriter, req *http.Request) {
	body := &webHookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		log.Println("Could not decode the request body ", err)
		return
	}

	parts := strings.Fields(body.Message.Text)

	helpText := "Supported commands:\n/english word - Define word with British English Dictionary\n/urban word " +
		"- Define word with Urban Dictionary\n/help - Display this help text"

	log.Println(len(parts))

	if len(parts) > 2 {
		if err := respond(body.Message.Chat.ID, "Please check your message and resend"); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
		return
	}
	log.Println(len(parts))
	command := parts[0]
	word := parts[1]

	if command == "/start" {
		welcomeMessage := "Hi.\nWelcome to Worder.\nTo get a definition, send an english word(en-gb) without unnecessary punctuations"
		if err := respond(body.Message.Chat.ID, welcomeMessage); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
	} else if command == "/help" {
		if err := respond(body.Message.Chat.ID, helpText); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
	} else if command == "/english" {
		definition := getDefinition(word)
		if err := respond(body.Message.Chat.ID, definition); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
	} else if command == "/urban" {
		definition := getUrbanDefinition(word)
		if err := respond(body.Message.Chat.ID, definition); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
	} else {
		if err := respond(body.Message.Chat.ID, "Unsupported format\n\n" + helpText); err != nil {
			log.Println("Error in sending message ", err)
			return
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
