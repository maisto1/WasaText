package schemas

// Represent a single comment
type Comment struct {
	ID             int64  `json:"id" example:"1"`                           // IUNique identifier
	CommentContent string `json:"comment-content" example:"Nice photo!"`    // Comment's content
	Username       string `json:"username" example:"Alessandro"`            // Comment's author
	Timestamp      string `json:"timestamp" example:"2024-11-15T14:30:00Z"` // Comment's timestamp
}
