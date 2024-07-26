package assistant

import (
	"eng_bot/knowledge/chat_gpt"
	"strings"
)

const (
	Api_key_set_response = "Api key has been set"
	Promt_set_response   = "Promt has been set"
)

type Assistant struct {
	UserId             string
	chatGptApiKey      string
	defaultSystemPromt string
	userSystemPromt    string
	chat_gpt           *chat_gpt.ChatGpt
}

var assisstantStorage = make(map[string]*Assistant)

func Call(user_id string) *Assistant {
	// check if assistant already exists
	if assistant, ok := assisstantStorage[user_id]; ok {
		return assistant
	}
	// create new assistant
	assistant := Assistant{
		UserId:             user_id,
		chatGptApiKey:      "",
		defaultSystemPromt: "Hello, You are my adult friend. But also you are a teacher of English. So if i'm wrong, please correct me,I'm here to learn and then response like a friend to my frase or question",
		userSystemPromt:    "",
		chat_gpt:           chat_gpt.New(""),
	}

	assisstantStorage[user_id] = &assistant
	return &assistant
}

func (a *Assistant) Ask(message string) string {

	command_result, is_command := a.handleMessage(message)

	if is_command {
		return command_result
	}
	if strings.TrimSpace(message) == "" {
		return "Empty message"
	}

	if a.chatGptApiKey == "" {
		return "Api key is not set"
	}

	return a.chat_gpt.Ask(a.chatGptApiKey, a.defaultSystemPromt, message)
}

func (a *Assistant) handleMessage(message string) (string, bool) {

	if message[0] == '/' {
		command_data := strings.Split(message, "=")
		if len(command_data) == 2 {
			command := command_data[0]
			command_value := command_data[1]
			if command == "/set_api_key" {
				a.chatGptApiKey = command_value
				return Api_key_set_response, true
			}

			if command == "/set_promt" {
				a.userSystemPromt = command_value
				return Promt_set_response, true
			}
		}

	}

	return message, false
}
