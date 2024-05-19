package main

import (
	"chatGptBot/telegram"
	"fmt"
	"log"
	"os"

	"chatGptBot/knowledge/chat_gpt"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("program start")
	token_telegram, token_gpt := token_must()
	telegram_client := telegram.New(token_telegram)
	chat_gpt := chat_gpt.New(token_gpt)
	//telegram_client.Start(chat_gpt)
	telegram_client.StartWebHook(chat_gpt, token_telegram)

}
func token_must() (string, string) {

	var token_tm string
	var token_gpt string
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token_tm = os.Getenv("TOKEN_TG")
	token_gpt = os.Getenv("TOKEN_GPT")
	if token_gpt == "" || token_tm == "" {
		log.Fatal("Empty tg or gpt token")
	}
	return token_tm, token_gpt
}
