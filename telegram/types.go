package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Update_Id int              `json:"update_id"`
	Message   *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}

type Chat struct {
	Id int `json:"id"`
}
