package models

type Preview struct {
	Conversation_id int64   `json:"id"`
	Name            string  `json:"name"`
	Photo           []byte  `json:"conversationPhoto"`
	LatestMessage   Message `json:"latestMessage"`
}
