package models

type Comment struct {
	Comment_id int64  `json:"id"`
	Message_id int64  `json:"message_id"`
	Content    string `json:"content"`
	Sender     User   `json:"sender"`
	Timestamp  int64  `json:"timestamp"`
}
