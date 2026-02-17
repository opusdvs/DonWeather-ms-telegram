package domain

type Update struct {
	Message Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	From From   `json:"from"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type From struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
