package models

// User represents a single user in the system.
type User struct {
	ID            int            `json:"id"`
	Username      string         `json:"username"`
	ProfilePhoto  string         `json:"profile-photo"`
	Conversations []Conversation `json:"conversations"`
}

// Conversation represents a conversation the user is part of.
type Conversation struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	ConversationPhoto string  `json:"conversation-photo"`
	LastMessage       Message `json:"last-message"`
}

// Message represents a message in the conversation.
type Message struct {
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
