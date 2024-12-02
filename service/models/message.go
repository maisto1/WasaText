package models

type Message struct {
	Message_id int64  `json:"id"`
	Timestamp  int64  `json:"timestamp"`
	Sender     User   `json:"sender"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Media      []byte `json:"media"`
	Status     string `json:"status"`
	Forwarded  bool   `json:"isForwarded"`
}
