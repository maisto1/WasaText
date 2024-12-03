package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	//Login or Register route
	rt.router.POST("/session", rt.wrap(rt.Login))

	//Get conversations preview
	rt.router.GET("/conversations", rt.wrap(rt.GetPreviewConversations))

	//Get users infos
	rt.router.GET("/users", rt.wrap(rt.GetUsers))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
