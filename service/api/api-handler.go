package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Login or Register route
	rt.router.POST("/session", rt.wrap(rt.Login, false))

	// Get conversations preview
	rt.router.GET("/conversations/", rt.wrap(rt.GetPreviewConversations, true))

	// Create a new conversation (group or private)
	rt.router.POST("/conversations/", rt.wrap(rt.CreateConversation, true))

	// Get messages of a specific conversation
	rt.router.GET("/conversations/:ConversationId", rt.wrap(rt.GetMessages, true))

	// Send a message in a specific conversation
	rt.router.POST("/conversations/:ConversationId/messages/", rt.wrap(rt.CreateMessage, true))

	// Delete a message from a conversation
	rt.router.DELETE("/conversations/:ConversationId/messages/:MessageId", rt.wrap(rt.DeleteMessage, true))

	// Forward a message to another conversation
	rt.router.POST("/conversations/:ConversationId/messages/:MessageId", rt.wrap(rt.ForwardMessage, true))

	// Reply to a message
	rt.router.POST("/conversations/:ConversationId/messages/:MessageId/reply", rt.wrap(rt.ReplyMessage, true))

	// Get a message's comment
	rt.router.GET("/conversations/:ConversationId/messages/:MessageId/comments/", rt.wrap(rt.GetComments, true))

	// Create a message's comment
	rt.router.POST("/conversations/:ConversationId/messages/:MessageId/comments/", rt.wrap(rt.CreateComment, true))

	// Delete a message from a conversation
	rt.router.DELETE("/conversations/:ConversationId/messages/:MessageId/comments/:CommentId", rt.wrap(rt.DeleteComment, true))

	// Add user to groupchat
	rt.router.POST("/conversations/:ConversationId/members/", rt.wrap(rt.AddGRoup, true))

	// Get all the group members
	rt.router.GET("/conversations/:ConversationId/members/", rt.wrap(rt.GetGroupMembers, true))

	// Remove or left groupchat
	rt.router.DELETE("/conversations/:ConversationId/members/:UserId", rt.wrap(rt.RemoveGRoup, true))

	// Update group name
	rt.router.PUT("/conversations/:ConversationId/name", rt.wrap(rt.EditName, true))

	// Update group photo
	rt.router.PUT("/conversations/:ConversationId/photo", rt.wrap(rt.EditPhoto, true))

	// Update profile name
	rt.router.PUT("/users/profile/username", rt.wrap(rt.EditProfileName, true))

	// Update profile photo
	rt.router.PUT("/users/profile/photo", rt.wrap(rt.EditProfilePhoto, true))

	// Get users infos
	rt.router.GET("/users/", rt.wrap(rt.GetUsers, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
