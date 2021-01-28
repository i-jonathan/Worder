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

	go processRequest(body)

}

// process request
func processRequest(update *webHookReqBody) {
	parts := strings.Fields(update.Message.Text)

	helpText := "Supported commands:\n/english word - Define word with British English Dictionary\n/urban word " +
		"- Define word with Urban Dictionary\n/help - Display this help text"

	log.Println(len(parts))

	if len(parts) > 2 {
		if err := respond(update.Message.Chat.ID, "Please check your message and resend"); err != nil {
			log.Println("Error in sending message ", err)
			return
		}
	}

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
		if err := respond(update.Message.Chat.ID, helpText); err != nil {
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

type spellCheck struct {
	Matches	[]matches	`json:"matches"`
}

type matches struct {
	Replacements	[]replacement	`json:"replacements"`
}

type replacement struct {
	Value	string	`json:"value"`
}

/* func grammarChecker(word, lang string) string {
	apiUrl := "https://grammarbot.p.rapidapi.com/check"
	payload :=strings.NewReader(fmt.Sprintf("text=%s&language=%s", word, lang))

	client := &http.Client{}
	grammar := &spellCheck{}

	req, err := http.NewRequest("POST", apiUrl, payload)
	if err != nil {
		log.Println("Grammar checker issue ", err)
		return "Check your spelling and retry."
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-rapidapi-key", os.Getenv("grammarKey"))
	req.Header.Set("x-rapidapi-host", os.Getenv("grammarHost"))

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("No suggestions: ", err)
		return "No suggestions."
	}

	err = json.NewDecoder(resp.Body).Decode(grammar)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Json decoder issue in grammar check: ", err)
		return "No Suggestions Found"
	}

	limit := 7
	count := 0
	suggestions := "Suggestions:"

	if len(grammar.Matches[0].Replacements) < limit {
		limit = len(grammar.Matches[0].Replacements)
	}

	for count < limit {
		suggestions += grammar.Matches[0].Replacements[count].Value + ","
	}

	return suggestions
} */