package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	//Login or Register route
	rt.router.POST("/session", rt.wrap(rt.Login, false))

	//Get conversations preview
	rt.router.GET("/conversations", rt.wrap(rt.GetPreviewConversations, true))

	rt.router.POST("/conversations", rt.wrap(rt.CreateConversation, true))

	//Get users infos
	rt.router.GET("/users", rt.wrap(rt.GetUsers, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
