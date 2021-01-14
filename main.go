package main

import (
	"log"
	"net/http"
	"os"
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
