package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
	"github.com/maisto1/WasaText/service/models"
)

func (rt *_router) GetMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Get Message: "
	var messages []models.Message

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "invalid conversation_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err = rt.db.GetMessages(ctx.User_id, conversation_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "conversation not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "messages sended to client")
}
