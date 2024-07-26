package main

import (
	"fmt"
	"log"
	"os"

	"eng_bot/telegram"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("program start")
	token_telegram := token_must()
	telegram_client := telegram.New(token_telegram)
	telegram_client.StartWebHook(token_telegram)

}

func token_must() string {

	var token_tm string
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token_tm = os.Getenv("TOKEN_TG")

	if token_tm == "" {
		log.Fatal("Empty tg or gpt token")
	}
	return token_tm
}
