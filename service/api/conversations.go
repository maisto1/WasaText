package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
	"github.com/maisto1/WasaText/service/models"
)

func (rt *_router) GetPreviewConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Preview Conversation: "
	var previews []models.Preview

	previews, err := rt.db.GetPreviewConversations(ctx.User_id)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(previews)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message + "conversation preview sended to client")

}

func isValidConversation(conversationType, groupName, partecipant string) bool {
	switch conversationType {
	case "private":
		return groupName == "" && partecipant != ""
	case "group":
		return groupName != "" && partecipant == ""
	default:
		return false
	}
}

func (rt *_router) CreateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Create Conversation: "

	var requestBody struct {
		GroupName   string `json:"groupName"`
		ConvType    string `json:"conversationType"`
		Partecipant string `json:"partecipant"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&requestBody)
	if err != nil || !(isValidConversation(requestBody.ConvType, requestBody.GroupName, requestBody.Partecipant)) {
		ctx.Logger.WithError(err).Error(message + "error decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.CreateConversation(ctx.User_id, requestBody.GroupName, requestBody.ConvType, requestBody.Partecipant)
	if err != nil {
		if err.Error() == "partecipant not found" {
			rt.baseLogger.WithError(err).Error(message + "partecipant not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err.Error() == "already have a conversation" {
			rt.baseLogger.WithError(err).Error(message + "user already have a conversation with this partecipant")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			rt.baseLogger.WithError(err).Error(message + "error checking if user exists")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	message = message + "conversation created\n"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	ctx.Logger.Info(message)

}
