package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
)

func (rt *_router) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Login: "

	var requestBody struct {
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(requestBody.Username) < 3 || len(requestBody.Username) > 16 {
		ctx.Logger.WithError(err).Error(message + "username doesn't respect pattern")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := rt.db.Login(requestBody.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Error(message + "error checking if user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message = message + "Authenticated user " + requestBody.Username + " with ID " + strconv.FormatInt(userId, 10) + "\n"

	response := map[string]int64{"id": userId}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message)
}
