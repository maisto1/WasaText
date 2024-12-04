package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
)

func (rt *_router) AddGRoup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Add Group: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "invalid conversation_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.AddGroup(ctx.User_id, requestBody.Username, conversation_id)
	if err != nil {
		if err.Error() == "user doesn't have private conversation" || err.Error() == "this isn't a group chat" {
			ctx.Logger.WithError(err).Error(message + "conversation  not valid or no private conversation")
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			ctx.Logger.WithError(err).Error(message + "conversation  not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Info(message + "user successfully added")
}

func (rt *_router) RemoveGRoup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Add Group: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "invalid conversation_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user_id_str := ps.ByName("UserId")
	user_id, err := strconv.ParseInt(user_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "invalid user_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.RemoveGroup(ctx.User_id, conversation_id, user_id)
	if err != nil {
		if err.Error() == "this isn't a group chat" {
			ctx.Logger.WithError(err).Error(message + "conversation  not valid or no private conversation")
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			ctx.Logger.WithError(err).Error(message + "conversation  not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Info(message + "user successfully removed")
}
