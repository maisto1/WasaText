package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
)

func (rt *_router) GetConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Conversations"

	userIdParam := ps.ByName("userId")

	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + " invalid userId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Getting the user's conversations
	response, err := rt.db.GetConversations(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + " user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Error in serializing json
		ctx.Logger.WithError(err).Error(message + " error in serializing json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + " sended to user: " + strconv.FormatInt(userId, 10))
}
