package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
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
		ctx.Logger.WithError(err).Error(message + "error parsing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message)
}
