package main

import (
	"dingtalk"
	"net/http"
)

func main() {
	println("starting --------")
	http.Handle("/gitlab/webhook", &dingtalk.GitlabWebhookHandler{})
	_ = http.ListenAndServe(":8000", nil)
}
