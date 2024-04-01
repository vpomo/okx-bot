package app

import (
	"net/http"
	util "okx-bot/restservice/utils"
)

var NotFoundHandler = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		util.Respond(w, util.Message(false, "This resources was not found on our server"))
		next.ServeHTTP(w, r)
	})
}
