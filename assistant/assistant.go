package assistant

import (
	"eng_bot/knowledge/chat_gpt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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
	history            []string
}

func (a *Assistant) getContextPromt() string {
	//cobine them into one string
	history_string := ""
	if len(a.history) == 0 {
		return ""
	}
	for _, message := range a.history {
		history_string += message + "\n"
	}
	return "Here you are the history of messaging from first to last: " + history_string
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

	promt := a.defaultSystemPromt + a.getContextPromt()

	answer := a.chat_gpt.Ask(a.chatGptApiKey, promt, message)

	a.addHistory(message, answer)

	return answer
}

func (a *Assistant) addHistory(message string, answer string) {
	// check if history length is more than 5 remove the oldest message
	// history_limit get from env HISTORY_LIMIT
	godotenv.Load()
	history_limit := os.Getenv("HISTORY_LIMIT")
	limit, err := strconv.Atoi(history_limit)
	if err != nil {
		limit = 10
	}
	if len(a.history) > limit {
		a.history = a.history[1:]
	}
	history_message := "User: " + message + "; You: " + answer
	a.history = append(a.history, history_message)
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
