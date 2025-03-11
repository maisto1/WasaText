package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
	"github.com/maisto1/WasaText/service/constants"
	"github.com/maisto1/WasaText/service/models"
)

func (rt *_router) GetMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Get Messages: "
	var messages []models.Message

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
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
		ctx.Logger.WithError(err).Error(message + constants.ErrParsing)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "messages sended to client")
}

func isValidMessage(mediaType, content string, media []byte) bool {
	switch mediaType {
	case "text":
		return len(content) > 0 && len(media) == 0
	case "media":
		return len(media) > 0
	default:
		return false
	}
}

func (rt *_router) CreateMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Create Message: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		Media   []byte `json:"media"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&requestBody)
	if err != nil || !(isValidMessage(requestBody.Type, requestBody.Content, requestBody.Media)) {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mess, err := rt.db.CreateMessage(ctx.User_id, conversation_id, 0, requestBody.Type, requestBody.Content, requestBody.Media, false)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "conversation not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(mess)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrParsing)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "messages sended to client")
}

func (rt *_router) DeleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Delete Message: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message_id_str := ps.ByName("MessageId")
	message_id, err := strconv.ParseInt(message_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "invalid message_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(ctx.User_id, conversation_id, message_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "user/conversation/message not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message + "message deleted successfully")
}

func (rt *_router) ForwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Forward Message: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message_id_str := ps.ByName("MessageId")
	message_id, err := strconv.ParseInt(message_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestBody struct {
		TargetConversationId int64 `json:"conversationId"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mess, err := rt.db.ForwardMessage(ctx.User_id, conversation_id, requestBody.TargetConversationId, message_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "user/conversation/message not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(mess)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrParsing)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "messages forwarded to another conversation")
}

func (rt *_router) ReplyMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Reply Message: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message_id_str := ps.ByName("MessageId")
	message_id, err := strconv.ParseInt(message_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "invalid message_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		Media   []byte `json:"media"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&requestBody)
	if err != nil || !(isValidMessage(requestBody.Type, requestBody.Content, requestBody.Media)) {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	replyMessage, err := rt.db.ReplyToMessage(
		ctx.User_id,
		conversation_id,
		message_id,
		requestBody.Type,
		requestBody.Content,
		requestBody.Media,
	)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "failed to create reply message")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(replyMessage)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrParsing)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "reply message sent successfully")
}
