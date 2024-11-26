package schemas

// Represent a single user
type User struct {
	ID            int        `json:"id" example:"1"`                                                  // Unique  identifier
	Username      string     `json:"username" example:"Alessandro"`                                   // Username
	ProfilePhoto  string     `json:"profile_photo,omitempty" example:"https://example.com/photo.jpg"` // URL profile pic
	Conversations []UserConv `json:"conversations,omitempty"`                                         // Conversation list
}

// Represent a single user's conversation
type UserConv struct {
	ID                int    `json:"id" example:"1"`                                             // Unique  identifier
	Name              string `json:"name" example:"Alessandro"`                                  // Cnversation name
	ConversationPhoto string `json:"conversation_photo" example:"https://example.com/photo.jpg"` // Conversation photo
	LastMessage       *struct {
		Content   string `json:"content" example:"Goodbye!"`               // LAst message
		Timestamp string `json:"timestamp" example:"2023-11-15T14:28:00Z"` // Timestamp
	} `json:"last_message,omitempty"`
}
