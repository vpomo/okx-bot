package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"okx-bot/restservice/app"
	"okx-bot/restservice/controllers"
	"os"
	"os/signal"
	"syscall"
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

	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	//router.NotFoundHandler = http.NotFoundHandler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	logger.Info("port: ", port)

	go func() {
		logger.Infoln("Serving REST started")
		err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
		if err != nil {
			logger.Error(err)
		}
	}()

	<-c
	logger.Info("Server graceful stopped")
	os.Exit(0)
}
