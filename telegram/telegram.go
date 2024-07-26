package telegram

import (
	"bytes"
	"encoding/json"
	"eng_bot/assistant"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	TELEGRAM_API      = "api.telegram.org"
	BATCH_SIZE        = 100
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type TelegramClient struct {
	host        string
	api_path    string
	http_client http.Client
}

type sendMessageReqBody struct {
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type helper interface {
	Ask(message string) string
}

func New(token string) *TelegramClient {
	return &TelegramClient{
		host:        TELEGRAM_API,
		api_path:    api_path(token),
		http_client: http.Client{},
	}
}

func api_path(token string) string {
	return "bot" + token
}

func webhookHandler(c *gin.Context, t *TelegramClient) {
	defer c.Request.Body.Close()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}
	assistant := getAssistant(update)

	question := update.Message.Text
	answer := assistant.Ask(question)
	t.Answer(answer, update.Message.Chat.Id)
}

func getAssistant(messageUpdate Update) *assistant.Assistant {
	telegram_user_id := strconv.Itoa(messageUpdate.Message.From.Id)
	
	return assistant.Call(telegram_user_id)
}

func (t *TelegramClient) StartWebHook(t_token string) {
	port := os.Getenv("LSTEN_PORT")
	hook_url := os.Getenv("HOOK_URL")
	if port == "" || hook_url == "" {
		log.Fatal("Empty listen port or hook url in env file")
	}
	router := gin.New()
	router.POST("/"+hook_url, func(c *gin.Context) {
		webhookHandler(c, t)
	})
	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}

func (t *TelegramClient) Messages(offset int, limit int) []Update {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	u := url.URL{
		Scheme: "https",
		Host:   t.host,
		Path:   path.Join(t.api_path, getUpdatesMethod),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil
	}

	req.URL.RawQuery = q.Encode()

	resp, err := t.http_client.Do(req)

	if err != nil {
		return nil
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	var response UpdatesResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	return response.Result
}

func (t *TelegramClient) Answer(answer string, chat_id int) {
	//send msg to telegram
	reqBody := sendMessageReqBody{
		ChatID: chat_id,
		Text:   answer,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Ошибка при маршализации запроса:", err)
		return
	}

	u := url.URL{
		Scheme: "https",
		Host:   t.host,
		Path:   path.Join(t.api_path, sendMessageMethod),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqBytes))

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
	}

	resp, err := t.http_client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return
	}

	fmt.Println("Ответ от сервера:", string(body))

}

func validateUser(messageUpdate Update) bool {

	config_allow_user := os.Getenv("ALLOW_USERS")
	if config_allow_user == "" {
		return false
	}
	allowed_user_ids := strings.Split(config_allow_user, ",")
	if len(allowed_user_ids) > 0 {
		message_user_id := strconv.Itoa(messageUpdate.Message.From.Id)
		for _, allowed_user_id := range allowed_user_ids {
			if allowed_user_id == message_user_id {
				return true
			}
		}
	}

	return false
}
