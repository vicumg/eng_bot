package chat_gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

const (
	GPT_API = "api.openai.com"
)

type ChatGpt struct {
	gpt_api   string
	gpt_token string
}

type gpt_message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var promt string = ""
var promt_in_text bool = false

type Gtp_request struct {
	Model    string        `json:"model"`
	Messages []gpt_message `json:"messages"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Choice struct {
	FinishReason string   `json:"finish_reason"`
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	Logprobs     []string `json:"logprobs"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ServerResponse struct {
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}

func New(token string) *ChatGpt {
	return &ChatGpt{
		gpt_api:   GPT_API,
		gpt_token: token,
	}
}
func (gpt *ChatGpt) Ask(message string) string {

	response := requestGpt(message)

	return response
}

func requestGpt(text string) string {

	if is_command(text) {
		return run_command(text)
	}

	gpt_message_system := gpt_message{
		Role:    "system",
		Content: get_prompt(text),
	}
	//remove prompt from text if it exists
	if promt_in_text {
		text = text[:len(text)-len(promt)-2]
	}

	gpt_message_user := gpt_message{
		Role:    "user",
		Content: text,
	}

	model := get_model()

	request := Gtp_request{
		Model: model,
	}
	request.Messages = append(request.Messages, gpt_message_system)
	request.Messages = append(request.Messages, gpt_message_user)

	reqBytes, err := json.Marshal(request)

	if err != nil {
		return "error with gtp create request"
	}

	u := url.URL{
		Scheme: "https",
		Host:   GPT_API,
		Path:   path.Join("v1/chat/completions"),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqBytes))

	gpt_token := os.Getenv("TOKEN_GPT")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gpt_token)
	if err != nil {
		fmt.Println("Error with request create:", err)
	}
	http_client := http.Client{}
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Error with request send:", err)
		return "Error with request send"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while answer reading:", err)
		return "Error while answer reading:"
	}

	var result ServerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can't unmarshal JSON")
	}

	return result.Choices[0].Message.Content
}

func get_model() string {
	config_file := "config.txt"
	file, err := os.Open(config_file)
	if err != nil {
		log.Fatal("Failed to open config file:", err)
	}
	defer file.Close()

	var model string
	_, err = fmt.Fscanf(file, "model_name:%q", &model)
	if err != nil {
		log.Fatal("Failed to read model name from config file:", err)
	}
	return model
}

func get_prompt(text string) string {

	start := 0
	end := 0
	for i, c := range text {
		if c == '[' {
			start = i
		}
		if c == ']' {
			end = i
		}
	}
	if start == 0 && end == 0 {
		promt_in_text = false
		return promt
	}
	promt_in_text = true
	promt = text[start+1 : end]
	return promt
}

func is_command(text string) bool {
	if text[0] == '/' {
		return true
	}
	return false
}

func run_command(text string) string {
	switch text {
	case "/start":
		return "Hello, I'm a chatbot. I can talk to you"
	case "/help":
		return "I can talk to you. Just write me a message"
	case "/get_promt":
		return promt
	default:
		return "I don't know this command"
	}
}