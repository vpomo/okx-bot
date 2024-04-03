package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/restservice/models"
	u "okx-bot/restservice/utils"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.signal-controllers",
	})
)

var ReceiveSignal = func(w http.ResponseWriter, r *http.Request) {

	signal := &models.TradingViewSignal{}

	err := json.NewDecoder(r.Body).Decode(signal)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	signal.Save()
	u.Respond(w, u.Message(true, "The signal is received"))
}
