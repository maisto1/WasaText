package schemas

// Represent a single message
type Message struct {
	ID             int       `json:"id" example:"1"`                                          // Unique identifier
	Sender         string    `json:"sender" example:"Alessandro"`                             // Sender
	Timestamp      string    `json:"timestamp" example:"2023-11-15T14:28:00Z"`                // Timestamp
	MediaType      string    `json:"mediatype" example:"text"`                                // Media type of message
	MessageContent string    `json:"message-content,omitempty" example:"Hi!"`                 // Message content
	Photo          string    `json:"photo,omitempty" example:"https://example.com/photo.jpg"` // Embedded message photo
	Status         string    `json:"status" example:"sent"`                                   // Message status
	IsForwarded    bool      `json:"is-forwarded" example:"false"`                            // Check if message is forwarded
	Comments       []Comment `json:"comments,omitempty"`                                      // Comment list
}
