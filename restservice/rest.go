package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/restservice/app"
	"okx-bot/restservice/controllers"
	"okx-bot/restservice/models"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logger = logrus.WithFields(logrus.Fields{
		"app":       "okx-bot",
		"component": "app.main-rest",
	})
)

func main() {

	logger.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	router.HandleFunc("/api/signal/receive", controllers.ReceiveSignal).Methods("POST") //  user/2/contacts

	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	//router.NotFoundHandler = http.NotFoundHandler()

	go func() {
		logger.Infoln("Waiting 5 second ...")
		time.Sleep(5 * time.Second)
		models.ConnectDB()
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	go func() {
		logger.Infoln("Server REST starting ...")
		err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
		if err != nil {
			logger.Error(err)
		}
		logger.Infoln("Serving REST started")
	}()

	<-c
	logger.Info("Server graceful stopped")
	os.Exit(0)
}
