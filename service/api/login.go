package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/maisto1/WasaText/service/api/reqcontext"
)

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	message := "Login:\n"

	var requestBody struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user exists in the database, if it does not exist, create it
	userID, err := rt.db.Login(requestBody.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Error(message + "error checking if user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	message = message + "Authenticated user " + requestBody.Username + " with ID " + strconv.FormatUint(userID, 10) + "\n"

	// response contains the user ID associated to the username
	response := map[string]uint64{"id": userID}

	// Send the response to the client
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error(message + "error parsing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(message+"sent UserID ", response["userId"], "to the client")
}
