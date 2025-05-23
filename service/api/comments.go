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

func (rt *_router) GetComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Get Comments: "
	var comments []models.Comment

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

	comments, err = rt.db.GetComments(ctx.User_id, conversation_id, message_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "conversation or message not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrParsing)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "comments sended to client")
}

func (rt *_router) CreateComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Create comment: "

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
		Content string `json:"content"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := rt.db.CreateComment(ctx.User_id, conversation_id, message_id, requestBody.Content)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "conversation or message not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrParsing)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "comment sended to client")

}

func (rt *_router) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Delete Message: "

	conversation_id_str := ps.ByName("ConversationId")
	conversation_id, err := strconv.ParseInt(conversation_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.InvalidConvId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment_id_str := ps.ByName("CommentId")
	comment_id, err := strconv.ParseInt(comment_id_str, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "invalid comment_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteComment(ctx.User_id, conversation_id, comment_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "user/conversation/message not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message + "comment deleted successfully")
}
