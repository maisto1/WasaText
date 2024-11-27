package api

import (
	"encoding/json"
	"net/http"
	"regexp"
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

	if len(requestBody.Username) < 2 || len(requestBody.Username) > 17 {
		ctx.Logger.WithError(err).Error(message + "username doesn't match length")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9_]*$`)
	if !re.MatchString(requestBody.Username) {
		ctx.Logger.WithError(err).Error(message + "username doesn't match pattern")
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
	message = message + "Authenticated username " + requestBody.Username + " with ID " + strconv.FormatInt(userID, 10) + "\n"

	// response contains the user ID associated to the username
	response := map[string]int64{"id": userID}

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
