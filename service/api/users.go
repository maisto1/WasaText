package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
	"github.com/maisto1/WasaText/service/constants"
	"github.com/maisto1/WasaText/service/models"
)

func (rt *_router) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Users:\n"
	var users []models.User

	usernames := r.URL.Query().Get("username")

	users = rt.db.GetUsers(usernames)

	message = message + "serching users -> " + usernames

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrParsing)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message)
}

func (rt *_router) EditProfilePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Edit Profile photo: "

	var requestBody struct {
		Photo []byte `json:"photo"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.EditProfilePhoto(ctx.User_id, requestBody.Photo)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "user  not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	ctx.Logger.Info(message + "photo updated successfully")
}

func (rt *_router) EditProfileName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Edit Profile name: "
	var requestBody struct {
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + constants.ErrDecBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.EditProfileName(ctx.User_id, requestBody.Username)
	if err != nil {
		if strings.Contains(err.Error(), "username already exists") {
			ctx.Logger.WithError(err).Error(message + "username already exists")
			w.WriteHeader(http.StatusConflict)
			return
		}

		ctx.Logger.WithError(err).Error(message + "user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info(message + "username updated successfully")
}
