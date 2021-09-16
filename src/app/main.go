package main

import (
	"dingtalk/src/dingtalk"
	"dingtalk/src/log"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	log.Logger.Formatter = new(logrus.JSONFormatter)
	log.Logger.Level = logrus.InfoLevel
	log.Logger.Out = os.Stdout
	file, err := os.OpenFile("/go/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Logger.Out = file
	} else {
		log.Logger.Info("Failed to log to file, using default stderr")
	}
}

func main() {
	log.Logger.Info("dingtalk is starting")
	http.Handle("/gitlab/webhook", &dingtalk.GitlabWebhookHandler{})
	_ = http.ListenAndServe(":8000", nil)
}
