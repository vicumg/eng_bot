package chat_gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	gpt_message_system := gpt_message{
		Role:    "system",
		Content: "You are software engeneer",
	}
	gpt_message_user := gpt_message{
		Role:    "user",
		Content: text,
	}

	request := Gtp_request{
		Model: "gpt-4-turbo",
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
		fmt.Println("Ошибка при создании запроса:", err)
	}
	http_client := http.Client{}
	resp, err := http_client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return "Ошибка при отправке запроса"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return "Ошибка при чтении тела ответа:"
	}

	fmt.Println("Ответ от сервера:", string(body))

	var result ServerResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return result.Choices[0].Message.Content
}
