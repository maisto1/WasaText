package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
	"github.com/maisto1/WasaText/service/constants"
)

func (rt *_router) AddGRoup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Add Group: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
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
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
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
	message := "Remove Group: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
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

func (rt *_router) EditPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Edit Photo Group: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody struct {
		GroupPhoto []byte `json:"groupPhoto"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.EditPhoto(conversation_id, requestBody.GroupPhoto)
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

	ctx.Logger.Info(message + "photo updated successfully")
}

func (rt *_router) EditName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Edit Name Group: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var requestBody struct {
		GroupName string `json:"groupName"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.EditName(conversation_id, requestBody.GroupName)
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

	ctx.Logger.Info(message + "name updated successfully")
}

func (rt *_router) GetGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Get Group Members: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	members, err := rt.db.GetGroupMembers(conversation_id)
	if err != nil {
		if err.Error() == "this isn't a group chat" {
			ctx.Logger.WithError(err).Error(message + "conversation not valid")
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			ctx.Logger.WithError(err).Error(message + "error retrieving group members")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(members)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "members retrieved successfully")
}
