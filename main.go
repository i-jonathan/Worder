package main

import (
	"log"
	"net/http"
	"os"
)

type webHookReqBody struct {
	Message message `json:"message"`
}

type message struct {
	Text string `json:"text"`
	Chat chat   `json:"chat"`
}

type chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type reply struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var token = os.Getenv("token")
var url = "https://api.telegram.org/bot" + token + "/"

func main() {
	//port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+"5000", http.HandlerFunc(handler))
	log.Println(err)
}
