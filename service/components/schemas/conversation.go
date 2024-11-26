package schemas

// Represent a single conversation
type Conversation struct {
	ID                int       `json:"id" example:"1"`                                             // Unique  identifier
	Name              string    `json:"name" example:"Alessandro"`                                  // Conversation name
	ConversationType  string    `json:"conversation-type" example:"text"`                           // Type of the conversation
	ConversationPhoto string    `json:"conversation-photo" example:"https://example.com/photo.jpg"` // Conversation photo
	Partecipants      []string  `json:"partecipants" example:"[\"Alessandro\"]"`                    // Conversation Partecipants
	Messages          []Message `json:"messages,omitempty"`                                         // Message list
}
